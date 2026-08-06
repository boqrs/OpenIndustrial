package role

import (
	"time"

	"github.com/google/uuid"
)

// Permission is a string that represents a specific action a user can perform.
// e.g., "users:create", "devices:read", "billing:manage"
type Permission string

// Role represents a collection of permissions that can be assigned to users.
// Roles are defined at the organization level.
type Role struct {
	ID          uuid.UUID    `json:"id"`
	OrgID       uuid.UUID    `json:"org_id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// NewRole creates a new Role entity.
func NewRole(orgID uuid.UUID, name string, permissions []Permission) (*Role, error) {
	if name == "" {
		return nil, ErrRoleNameRequired
	}

	now := time.Now().UTC()
	return &Role{
		ID:          uuid.New(),
		OrgID:       orgID,
		Name:        name,
		Permissions: permissions,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}