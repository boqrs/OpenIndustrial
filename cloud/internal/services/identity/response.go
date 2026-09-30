package identity

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// =====================================================
// Authentication Response
// =====================================================

// LoginResponse
//
// 登录成功返回
//
// access token:
//
//	用于API访问
//
// refresh token:
//
//	用于刷新access token
type LoginResponse struct {

	// JWT Access Token
	//
	AccessToken string `json:"access_token"`

	// Refresh Token
	//
	RefreshToken string `json:"refresh_token"`

	// Access Token过期时间
	//
	ExpiresIn int64 `json:"expires_in"`

	// Token类型
	//
	// Bearer
	//
	TokenType string `json:"token_type"`

	// 当前登录用户
	//
	User *model.User `json:"user"`
}

// UserResponse
//
// 对外返回用户信息
//
// 避免直接暴露model
type UserResponse struct {
	ID uint `json:"id"`

	UUID uuid.UUID `json:"uuid"`

	TenantID uint `json:"tenant_id"`

	Email string `json:"email"`

	Name string `json:"name"`

	UserType string `json:"user_type"`

	Status string `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}

func NewUserResponse(
	user *model.User,
) *UserResponse {

	if user == nil {
		return nil
	}

	return &UserResponse{

		ID: user.ID,

		UUID: user.UUID,

		TenantID: user.TenantID,

		Email: user.Email,

		Name: user.Name,

		UserType: user.UserType,

		Status: user.Status,

		CreatedAt: user.CreatedAt,
	}

}

// =====================================================
// JWT
// =====================================================

// JWTClaims
//
// access token / refresh token
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`

	TenantID uint `json:"tenant_id"`

	Email string `json:"email"`

	TokenType string `json:"token_type"`

	jwt.RegisteredClaims
}
