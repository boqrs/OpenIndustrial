package product

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the storage interface for products and their instances.
type Repository interface {
	// Product methods
	CreateProduct(ctx context.Context, product *Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error)
	ListProductsByOrg(ctx context.Context, orgID uuid.UUID) ([]*Product, error)

	// ProductInstance methods
	CreateInstance(ctx context.Context, instance *ProductInstance) error
	GetInstanceBySN(ctx context.Context, sn string) (*ProductInstance, error)
	GetInstanceByID(ctx context.Context, id uuid.UUID) (*ProductInstance, error)

	// LifecycleEvent methods
	AppendLifecycleEvent(ctx context.Context, event *LifecycleEvent) error
}