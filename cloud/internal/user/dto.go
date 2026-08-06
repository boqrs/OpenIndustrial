package user

import "github.com/google/uuid"

// CreateUserRequest defines the structure for a request to create a new user.
type CreateUserRequest struct {
	OrgID    string `json:"org_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone,omitempty"`
}

// UserResponse defines the structure for a standard user API response.
// It safely exposes user data without the password hash.
type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	OrgID     uuid.UUID  `json:"org_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone,omitempty"`
	Status    UserStatus `json:"status"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

// ToUserResponse converts a User entity to a UserResponse DTO.
func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		OrgID:     user.OrgID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}