package param

import (
	"encoding/json"

	"github.com/google/uuid"
)

// RegisterTenantRequest defines the parameters for registering a new tenant.
type RegisterTenantRequest struct {
	TenantName    string `json:"tenant_name" binding:"required"`
	AdminEmail    string `json:"admin_email" binding:"required,email"`
	AdminPassword string `json:"admin_password" binding:"required,min=8"`
}

// LoginRequest defines the parameters for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TenantID string `json:"tenant_id" binding:"required"`
}

// CreateUserRequest defines the parameters for creating a new user.
type CreateUserRequest struct {
	UserType string          `json:"user_type" binding:"required"`
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=8"`
	RoleName string          `json:"role_name" binding:"required"`
	Profile  json.RawMessage `json:"profile"`
}

// ListUsersRequest defines the parameters for listing users from the API layer.
type ListUsersRequest struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

// UpdateUserRequest defines the parameters for updating a user.
type UpdateUserRequest struct {
	Profile  json.RawMessage `json:"profile"`
	RoleName string          `json:"role_name"`
}

// AssignRoleToUserRequest defines the parameters for assigning a role to a user.
type AssignRoleToUserRequest struct {
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

type ListUsersRepoReq struct {
	Limit  int
	Offset int
}
