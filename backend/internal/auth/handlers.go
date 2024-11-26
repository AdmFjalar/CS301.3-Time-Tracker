package auth

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

type AuthHandlers struct {
	Authenticator auth.Authenticator
	Store         store.Storage
	Mailer        mailer.Client
	Config        config
	Logger        *zap.SugaredLogger
}

func (h *AuthHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Email: payload.Email,
		Role: store.Role{
			Name: "user",
		},
	}

	if err := user.Password.Set(payload.Password); err != nil {
		internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := h.Store.Users.CreateAndInvite(ctx, user, hashToken, h.Config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			badRequestResponse(w, r, err)
		default:
			internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	activationURL := fmt.Sprintf("%s/confirm/%s", h.Config.frontendURL, plainToken)

	isProdEnv := h.Config.env == "production"
	vars := struct {
		ActivationURL string
	}{
		ActivationURL: activationURL,
	}

	status, err := h.Mailer.Send(mailer.UserWelcomeTemplate, user.Email, vars, !isProdEnv)
	if err != nil {
		h.Logger.Errorw("error sending welcome email", "error", err)

		if err := h.Store.Users.Delete(ctx, user.ID); err != nil {
			h.Logger.Errorw("error deleting user", "error", err)
		}

		internalServerError(w, r, err)
		return
	}

	h.Logger.Infow("Email sent", "status code", status)

	if err := jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		internalServerError(w, r, err)
	}
}

func (h *AuthHandlers) RequestPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPasswordResetPayload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user, err := h.Store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	hash := sha256.Sum256([]byte(uuid.New().String()))
	hashToken := hex.EncodeToString(hash[:])

	err = h.Store.Users.RequestPasswordAndEmailReset(r.Context(), user, hashToken, h.Config.mail.exp)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	isProdEnv := h.Config.env == "production"
	vars := struct {
		ResetURL string
	}{
		ResetURL: fmt.Sprintf("%s/reset-password/%s", h.Config.frontendURL, hashToken),
	}

	status, err := h.Mailer.Send(mailer.PasswordResetTemplate, user.Email, vars, !isProdEnv)
	if err != nil {
		h.Logger.Errorw("error sending password reset email", "error", err)
		internalServerError(w, r, err)
		return
	}

	h.Logger.Infow("Email sent", "status code", status)

	if err := jsonResponse(w, http.StatusOK, nil); err != nil {
		internalServerError(w, r, err)
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (h *AuthHandlers) CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user, err := h.Store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			unauthorizedErrorResponse(w, r, err)
		default:
			internalServerError(w, r, err)
		}
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		unauthorizedErrorResponse(w, r, err)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(h.Config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": h.Config.auth.token.iss,
		"aud": h.Config.auth.token.iss,
	}

	token, err := h.Authenticator.GenerateToken(claims)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, token); err != nil {
		internalServerError(w, r, err)
	}
}

func (h *AuthHandlers) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var payload ChangePasswordPayload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)

	if err := user.Password.Compare(payload.OldPassword); err != nil {
		unauthorizedErrorResponse(w, r, err)
		return
	}

	if err := user.Password.Set(payload.NewPassword); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := h.Store.Users.ChangePassword(r.Context(), user); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusNoContent, ""); err != nil {
		internalServerError(w, r, err)
	}
}
