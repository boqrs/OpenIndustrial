package auth

import "github.com/golang-jwt/jwt/v4"

// Claims defines the structure of the JWT payload.
type Claims struct {
	jwt.RegisteredClaims
	Type    string `json:"type"` // e.g., "human", "gateway", "device"
	Session string `json:"session,omitempty"`
}

// TokenService defines the interface for generating and validating tokens.
type TokenService interface {
	Generate(claims Claims) (accessToken string, refreshToken string, err error)
	Validate(token string) (*Claims, error)
}