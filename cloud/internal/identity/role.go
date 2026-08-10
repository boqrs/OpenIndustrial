package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Role represents a role in the system.
type Role struct {
	ID          uuid.UUID `db:"id"`
	TenantID    uuid.UUID `db:"tenant_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	IsSystem    bool      `db:"is_system"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// RoleRepository defines the interface for role persistence.
type RoleRepository interface {
	CreateRole(ctx context.Context, role *Role) error
	GetRoleByID(ctx context.Context, tenantID, id uuid.UUID) (*Role, error)
	GetRoleByName(ctx context.Context, tenantID uuid.UUID, name string) (*Role, error)
	AddUserToRole(ctx context.Context, userID, roleID, tenantID uuid.UUID) error
}