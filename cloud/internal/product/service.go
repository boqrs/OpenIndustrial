package product

import (
	"context"

	"github.com/google/uuid"
)

// Service provides business logic for products.
type Service struct {
	repo Repository
}

// NewService creates a new product service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateProduct creates a new product type.
func (s *Service) CreateProduct(ctx context.Context, orgID uuid.UUID, name, code, spec, description string) (*Product, error) {
	product := NewProduct(orgID, name, code, spec, description)
	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

// ListProductsByOrg lists all products for a given organization.
func (s *Service) ListProductsByOrg(ctx context.Context, orgID uuid.UUID) ([]*Product, error) {
	return s.repo.ListProductsByOrg(ctx, orgID)
}

// GetProductByID retrieves a single product by its ID.
func (s *Service) GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

// AddLifecycleEvent adds a new lifecycle event to a product instance.
func (s *Service) AddLifecycleEvent(ctx context.Context, event *LifecycleEvent) error {
	// In a real system, you might add validation here, e.g., check if the product instance exists.
	return s.repo.AppendLifecycleEvent(ctx, event)
}