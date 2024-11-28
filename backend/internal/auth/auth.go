package auth

import "github.com/golang-jwt/jwt/v5"

// Authenticator is an interface for generating and validating JWT tokens.
type Authenticator interface {
	// GenerateToken generates a JWT token with the given claims.
	GenerateToken(claims jwt.Claims) (string, error)
	// ValidateToken validates the given JWT token and returns the parsed token.
	ValidateToken(token string) (*jwt.Token, error)
}
