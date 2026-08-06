package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// LoginRequest defines the structure for a login request.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse defines the structure for a successful login response.
type LoginResponse struct {
	Token string `json:"token"`
}