package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Role defines a set of permissions that can be assigned to users.
type Role struct {
	ID           uuid.UUID
	TenantID     uuid.UUID // Can be nil for system-wide roles
	Name         string
	Description  string
	IsSystemRole bool
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Relations
	Permissions []*Permission
}

// RoleRepository defines the interface for database operations on roles.
type RoleRepository interface {
	CreateRole(ctx context.Context, role *Role) error
	GetRoleByID(ctx context.Context, tenantID, roleID uuid.UUID) (*Role, error)
	GetRoleByName(ctx context.Context, tenantID uuid.UUID, name string) (*Role, error)
	AddUserToRole(ctx context.Context, userID, roleID, tenantID uuid.UUID) error
	RemoveUserFromRole(ctx context.Context, userID, roleID, tenantID uuid.UUID) error
	AddPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	ListRoles(ctx context.Context, tenantID uuid.UUID) ([]Role, error)
}