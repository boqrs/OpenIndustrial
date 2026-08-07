package role

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for role storage.
type Repository interface {
	CreateRole(ctx context.Context, role *Role) error
	GetRoleByID(ctx context.Context, roleID uuid.UUID) (*Role, error)
	ListRolesByOrg(ctx context.Context, orgID uuid.UUID) ([]*Role, error)
}