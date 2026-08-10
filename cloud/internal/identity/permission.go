package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Permission represents an action that can be performed on a resource.
type Permission struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"` // e.g., "users:create"
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// PermissionRepository defines the interface for permission-related database operations.
type PermissionRepository interface {
	// CheckPermissionForUser checks if a user has a specific permission through their roles.
	CheckPermissionForUser(ctx context.Context, userID uuid.UUID, permissionName string) (bool, error)

	// CreatePermission adds a new permission to the database.
	CreatePermission(ctx context.Context, p *Permission) error

	// GetPermission retrieves a permission by its key and action.
	GetPermission(ctx context.Context, resourceKey, action string) (*Permission, error)

	// ListPermissionsByRole retrieves all permissions associated with a specific role.
	ListPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*Permission, error)
}