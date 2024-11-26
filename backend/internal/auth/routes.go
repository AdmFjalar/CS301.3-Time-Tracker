package auth

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterAuthRoutes(r chi.Router, handlers *AuthHandlers) {
	r.Route("/authentication", func(r chi.Router) {
		r.Post("/user", handlers.RegisterUserHandler)
		r.Post("/token", handlers.CreateTokenHandler)
		r.Post("/request-password-reset", handlers.RequestPasswordResetHandler)
		r.Put("/reset-password/{token}", handlers.ResetPasswordHandler)
	})
}
