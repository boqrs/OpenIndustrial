package user

import (
	"time"
)

// CreateUserRequest defines the structure for a request to create a new user.
type CreateUserRequest struct {
	OrgID    string `json:"org_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest defines the structure for a user login request.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse defines the structure for a response containing user details.
// Importantly, it omits the PasswordHash.
type UserResponse struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToUserResponse converts a User entity to a UserResponse DTO.
func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:        user.ID.String(),
		OrgID:     user.OrgID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}