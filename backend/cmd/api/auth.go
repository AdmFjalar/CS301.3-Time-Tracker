package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/mailer"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

type ChangePasswordPayload struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=3,max=72"`
}

type ResetPasswordPayload struct {
	Password string `json:"password" validate:"required,min=3,max=72"`
	Email    string `json:"email" validate:"required,email,max=255"`
}

type RequestPasswordResetPayload struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

// registerUserHandler godoc
//
//	@Summary		Creates a user
//	@Description	Creates a user and sends a welcome email with an activation link
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User information"
//	@Success		201		{object}	UserWithToken			"User created"
//	@Failure		400		{object}	error
//	@Failure		409		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the request payload
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Initialize a new user object
	user := &store.User{
		Email: payload.Email,
		Role: store.Role{
			Name: "user",
		},
	}

	// Hash the user's password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	// Generate a plain text token for email confirmation
	plainToken := uuid.New().String()

	// Hash the token for storage
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// Store user and hashed token in the database
	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Prepare the response object
	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	// Create activation URL for the welcome email
	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	isProdEnv := app.config.env == "production"
	vars := struct {
		ActivationURL string
	}{
		ActivationURL: activationURL,
	}

	// Send the welcome email
	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// Rollback user creation if email fails (SAGA pattern)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent", "status code", status)

	// Send the response
	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}

// requestPasswordResetHandler godoc
//
//	@Summary		Requests a password reset
//	@Description	Sends a password reset link to the user's email
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RequestPasswordResetPayload	true	"Email"
//	@Success		200		{string}	string							"Password reset email sent"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/request-password-reset [post]
func (app *application) requestPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPasswordResetPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get the user by email
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(uuid.New().String()))
	hashToken := hex.EncodeToString(hash[:])

	// Store the token in the database with an expiration time (e.g., 1 hour)
	err = app.store.Users.RequestPasswordAndEmailReset(r.Context(), user, hashToken, app.config.mail.exp)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Construct reset URL
	resetURL := fmt.Sprintf("%s/reset-password/%s", app.config.frontendURL, hashToken)

	// Send mail
	isProdEnv := app.config.env == "production"
	vars := struct {
		ResetURL string
	}{
		ResetURL: resetURL,
	}

	status, err := app.mailer.Send(mailer.PasswordResetTemplate, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending password reset email", "error", err)
		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent", "status code", status)

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Read the payload from the request body
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the payload
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get the user by email
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorizedErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Compare the provided password to the user's password
	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	// Create the JWT token
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Send the token in the response
	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}
}

// changePasswordHandler godoc
//
//	@Summary		Change the user's password
//	@Description	Allows a user to change their password by providing the old and new passwords
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		ChangePasswordPayload	true	"Old and new passwords"
//	@Success		204		{string}	string	"No Content"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/change-password [put]
func (app *application) changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the request payload
	var payload ChangePasswordPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the payload structure
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Retrieve the user from the request context
	user := getUserFromContext(r)

	// Compare the old password with the stored password
	if err := user.Password.Compare(payload.OldPassword); err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	// Set the new password
	if err := user.Password.Set(payload.NewPassword); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Update the password in the database
	if err := app.store.Users.ChangePassword(r.Context(), user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Send a successful response with no content
	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}
