package identity

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`

	RefreshToken string `json:"refresh_token"`

	ExpiresIn int64 `json:"expires_in"`

	TokenType string `json:"token_type"`

	User *UserResponse `json:"user"`
}

type UserResponse struct {
	ID uint `json:"id"`

	UUID uuid.UUID `json:"uuid"`

	TenantID uint `json:"tenant_id"`

	RoleID uint `json:"role_id"`

	Email string `json:"email"`

	Name string `json:"name"`

	UserType string `json:"user_type"`

	Status string `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}

func NewUserResponse(user *model.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:        user.ID,
		UUID:      user.UUID,
		TenantID:  user.TenantID,
		RoleID:    user.RoleID,
		Email:     user.Email,
		Name:      user.Name,
		UserType:  user.UserType,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
}

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`

	TenantID uint `json:"tenant_id"`

	RoleID uint `json:"role_id"`

	TokenType string `json:"token_type"`

	jwt.RegisteredClaims
}
