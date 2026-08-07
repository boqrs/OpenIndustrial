package user

import (
	"time"
)

// CreateUserRequest defines the request for creating a user.
type CreateUserRequest struct {
	OrgID     string `json:"org_id" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserResponse is the DTO for a user.
type UserResponse struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// ToUserResponse converts a User entity to a DTO.
func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:        user.ID.String(),
		OrgID:     user.OrgID.String(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}
}