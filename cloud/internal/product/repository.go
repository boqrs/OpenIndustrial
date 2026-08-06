package product

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for data access operations for a Product.
type Repository interface {
	Create(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, orgID, productID uuid.UUID) (*Product, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, orgID, productID uuid.UUID) error
}