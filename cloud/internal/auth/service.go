package auth

import (
	"context"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	OrgID  uuid.UUID `json:"org_id"`
	jwt.RegisteredClaims
}

// Service provides authentication functionalities.
type Service struct {
	userRepo   user.Repository
	jwtSecret  []byte
	tokenExpiry time.Duration
}

// NewService creates a new auth service.
func NewService(userRepo user.Repository, jwtSecret string) *Service {
	return &Service{
		userRepo:   userRepo,
		jwtSecret:  []byte(jwtSecret),
		tokenExpiry: 24 * time.Hour, // Default to 24 hours
	}
}

// Authenticate a user by email and password and return a JWT.
func (s *Service) Authenticate(ctx context.Context, email, password string) (string, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if !u.CheckPassword(password) {
		return "", ErrInvalidCredentials
	}

	// Create JWT claims
	expirationTime := time.Now().Add(s.tokenExpiry)
	claims := &Claims{
		UserID: u.ID,
		OrgID:  u.OrgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and return it
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}