package role

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for data access operations for a Role.
type Repository interface {
	Create(ctx context.Context, role *Role) error
	FindByID(ctx context.Context, orgID, roleID uuid.UUID) (*Role, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]*Role, error)
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, orgID, roleID uuid.UUID) error
}