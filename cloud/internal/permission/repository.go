package permission

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for data access operations for a Permission.
type Repository interface {
	Create(ctx context.Context, p *Permission) error
	FindByID(ctx context.Context, permissionID uuid.UUID) (*Permission, error)
	FindByName(ctx context.Context, name string) (*Permission, error)
	ListAll(ctx context.Context) ([]*Permission, error)
}