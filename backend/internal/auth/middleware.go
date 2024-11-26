package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func AuthTokenMiddleware(authenticator Authenticator, getUser func(context.Context, int64) (*store.User, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			token := parts[1]
			jwtToken, err := authenticator.ValidateToken(token)
			if err != nil {
				unauthorizedErrorResponse(w, r, err)
				return
			}

			claims, _ := jwtToken.Claims.(jwt.MapClaims)

			userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
			if err != nil {
				unauthorizedErrorResponse(w, r, err)
				return
			}

			ctx := r.Context()

			user, err := getUser(ctx, userID)
			if err != nil {
				unauthorizedErrorResponse(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, userCtx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func BasicAuthMiddleware(username, pass string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				unauthorizedBasicErrorResponse(w, r, err)
				return
			}

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CheckTimestampOwnershipMiddleware(getUserFromContext func(*http.Request) *store.User, getTimestampFromCtx func(*http.Request) *store.Timestamp, checkRolePrecedence func(context.Context, *store.User, string) (bool, error)) func(string, http.HandlerFunc) http.HandlerFunc {
	return func(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := getUserFromContext(r)
			timestamp := getTimestampFromCtx(r)

			if timestamp.UserID == user.ID {
				next.ServeHTTP(w, r)
				return
			}

			allowed, err := checkRolePrecedence(r.Context(), user, requiredRole)
			if err != nil {
				internalServerError(w, r, err)
				return
			}

			if !allowed {
				forbiddenResponse(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CheckRolePrecedenceMiddleware(getUserFromContext func(*http.Request) *store.User, checkRolePrecedence func(context.Context, *store.User, string) (bool, error)) func(string, http.HandlerFunc) http.HandlerFunc {
	return func(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := getUserFromContext(r)

			allowed, err := checkRolePrecedence(r.Context(), user, requiredRole)
			if err != nil {
				internalServerError(w, r, err)
				return
			}

			if !allowed {
				forbiddenResponse(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RateLimiterMiddleware(rateLimiter ratelimiter.Limiter, enabled bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if enabled {
				if allow, retryAfter := rateLimiter.Allow(r.RemoteAddr); !allow {
					rateLimitExceededResponse(w, r, retryAfter.String())
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
