package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthenticator is a struct that holds the secret, audience, and issuer for JWT authentication.
type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

// NewJWTAuthenticator creates a new JWTAuthenticator with the given secret, audience, and issuer.
func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{secret, iss, aud}
}

// GenerateToken generates a JWT token with the given claims.
func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the given JWT token and returns the parsed token.
func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
