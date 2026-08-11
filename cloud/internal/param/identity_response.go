package param

import (
	"encoding/json"

	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

// RegisterTenantResponse defines the result of a tenant registration.
type RegisterTenantResponse struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	AdminUserID uuid.UUID `json:"admin_user_id"`
}

// LoginResponse defines the result of a successful login.
type LoginResponse struct {
	Token string `json:"token"`
}

// GetCurrentUserResponse defines the result for getting the current user.
type GetCurrentUserResponse struct {
	ID       uuid.UUID       `json:"id"`
	UserType string          `json:"user_type"`
	Profile  json.RawMessage `json:"profile"`
}

// CreateUserResponse defines the result of creating a new user.
type CreateUserResponse struct {
	ID uuid.UUID `json:"id"`
}

// UserResponse defines the public representation of a user.
type UserResponse struct {
	ID        uuid.UUID       `json:"id"`
	TenantID  uuid.UUID       `json:"tenant_id"`
	UserType  string          `json:"user_type"`
	Profile   json.RawMessage `json:"profile"`
	CreatedAt int64           `json:"created_at"`
}

// RoleResponse defines the public representation of a role.
type RoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsSystem    bool      `json:"is_system"`
}

// GroupResponse defines the public representation of a group.
type GroupResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// ToUserResponse converts a model.User to a public-facing UserResponse.
func ToUserResponse(user *model.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.UUID,
		TenantID:  user.TenantID,
		UserType:  user.UserType,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt.Unix(),
	}
}

// ToRoleResponse converts a model.Role to a public-facing RoleResponse.
func ToRoleResponse(role *model.Role) *RoleResponse {
	if role == nil {
		return nil
	}
	return &RoleResponse{
		ID:          role.UUID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
	}
}

// ToGroupResponse converts a model.Group to a public-facing GroupResponse.
func ToGroupResponse(group *model.Group) *GroupResponse {
	if group == nil {
		return nil
	}
	return &GroupResponse{
		ID:          group.UUID,
		Name:        group.Name,
		Description: group.Description,
	}
}