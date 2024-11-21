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

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Email: payload.Email,
		Role: store.Role{
			Name: "user",
		},
	}

	// hash the user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	plainToken := uuid.New().String()

	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

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

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	isProdEnv := app.config.env == "production"
	vars := struct {
		ActivationURL string
	}{
		ActivationURL: activationURL,
	}

	// send mail
	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// rollback user creation if email fails (SAGA pattern)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent", "status code", status)

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
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
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

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

	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

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

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var payload ChangePasswordPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)

	if err := user.Password.Compare(payload.OldPassword); err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	if err := user.Password.Set(payload.NewPassword); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Users.ChangePassword(r.Context(), user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}

}

// // PasswordResetRequestPayload struct for requesting password reset
// type PasswordResetRequestPayload struct {
// 	Email string `json:"email" validate:"required,email,max=255"`
// }

// // requestPasswordResetHandler godoc
// //
// //	@Summary		Requests a password reset
// //	@Description	Sends a password reset link to the user's email
// //	@Tags			authentication
// //	@Accept			json
// //	@Produce		json
// //	@Param			payload	body		PasswordResetRequestPayload	true	"Email"
// //	@Success		200		{string}	string							"Password reset email sent"
// //	@Failure		400		{object}	error
// //	@Failure		500		{object}	error
// //	@Router			/authentication/request-password-reset [post]
// func (app *application) requestPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload PasswordResetRequestPayload
// 	if err := readJSON(w, r, &payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := Validate.Struct(payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// Retrieve user by email
// 	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
// 	if err != nil {
// 		switch err {
// 		case store.ErrNotFound:
// 			// Don't reveal whether the email exists for security reasons
// 			app.jsonResponse(w, http.StatusOK, "If that email exists, a reset link will be sent.")
// 		default:
// 			app.internalServerError(w, r, err)
// 		}
// 		return
// 	}

// 	// Generate a reset token (UUID or JWT)
// 	resetToken := uuid.New().String()
// 	hash := sha256.Sum256([]byte(resetToken))
// 	hashToken := hex.EncodeToString(hash[:])

// 	// Store token in the database with an expiration time (e.g., 1 hour)
// 	err = app.store.Users.StorePasswordResetToken(r.Context(), user.ID, hashToken, time.Now().Add(time.Hour))
// 	if err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	// Construct reset URL
// 	resetURL := fmt.Sprintf("%s/reset-password/%s", app.config.frontendURL, resetToken)

// 	// Send reset email
// 	isProdEnv := app.config.env == "production"
// 	vars := struct {
// 		ResetURL string
// 	}{
// 		ResetURL: resetURL,
// 	}

// 	status, err := app.mailer.Send(mailer.PasswordResetTemplate, user.Email, vars, !isProdEnv)
// 	if err != nil {
// 		app.logger.Errorw("error sending password reset email", "error", err)
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	app.logger.Infow("Password reset email sent", "status code", status)

// 	app.jsonResponse(w, http.StatusOK, "If that email exists, a reset link will be sent.")
// }

// // ResetPasswordWithTokenPayload struct for resetting password with token
// type ResetPasswordWithTokenPayload struct {
// 	ResetToken string `json:"reset_token" validate:"required"`
// 	NewPassword string `json:"new_password" validate:"required,min=3,max=72"`
// }

// // resetPasswordWithTokenHandler godoc
// //
// //	@Summary		Resets the user's password with a reset token
// //	@Description	Allows a user to reset their password with a reset token
// //	@Tags			authentication
// //	@Accept			json
// //	@Produce		json
// //	@Param			payload	body		ResetPasswordWithTokenPayload	true	"Reset token and new password"
// //	@Success		200		{string}	string							"Password reset successfully"
// //	@Failure		400		{object}	error
// //	@Failure		401		{object}	error
// //	@Failure		500		{object}	error
// //	@Router			/authentication/reset-password [put]
// func (app *application) resetPasswordWithTokenHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload ResetPasswordWithTokenPayload
// 	if err := readJSON(w, r, &payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := Validate.Struct(payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// Look up the reset token
// 	user, err := app.store.Users.GetByResetToken(r.Context(), payload.ResetToken)
// 	if err != nil {
// 		switch err {
// 		case store.ErrNotFound:
// 			app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid or expired reset token"))
// 		default:
// 			app.internalServerError(w, r, err)
// 		}
// 		return
// 	}

// 	// Set the new password
// 	if err := user.Password.Set(payload.NewPassword); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	// Delete the reset token after successful password reset
// 	err = app.store.Users.DeleteResetToken(r.Context(), user.ID)
// 	if err != nil {
// 		app.logger.Errorw("error deleting reset token", "error", err)
// 	}

// 	// Update user in the database
// 	err = app.store.Users.Update(r.Context(), user)
// 	if err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	// Send success response
// 	if err := app.jsonResponse(w, http.StatusOK, "Password reset successfully"); err != nil {
// 		app.internalServerError(w, r, err)
// 	}
// }
