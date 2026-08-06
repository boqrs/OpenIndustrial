package product

import (
	"context"

	"github.com/google/uuid"
)

// Service encapsulates the business logic for the product domain.
type Service struct {
	repo Repository
}

// NewService creates a new product service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateProduct handles the business logic of creating a new product definition.
func (s *Service) CreateProduct(ctx context.Context, orgID uuid.UUID, name, description, model string) (*Product, error) {
	product, err := NewProduct(orgID, name, description, model)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product by its ID within a specific organization.
func (s *Service) GetProductByID(ctx context.Context, orgID, productID uuid.UUID) (*Product, error) {
	return s.repo.FindByID(ctx, orgID, productID)
}

// ListProductsForOrg lists all products for a given organization.
func (s *Service) ListProductsForOrg(ctx context.Context, orgID uuid.UUID) ([]*Product, error) {
	return s.repo.ListByOrg(ctx, orgID)
}