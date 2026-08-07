package product

import (
	"context"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/event"
)

// Service defines the business logic for the product domain.
type Service interface {
	CreateProduct(ctx context.Context, orgID, name, description string, spec map[string]string) (*Product, error)
	GetProductByID(ctx context.Context, orgID, productID string) (*Product, error)
	ListProductsForOrg(ctx context.Context, orgID string) ([]*Product, error)
	AppendLifecycle(ctx context.Context, event LifecycleEvent) error
}

// Repository defines the persistence interface for product related entities.
type Repository interface {
	Create(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, orgID, productID string) (*Product, error)
	ListByOrg(ctx context.Context, orgID string) ([]*Product, error)
	// Add a method to save lifecycle events
	SaveLifecycleEvent(ctx context.Context, event LifecycleEvent) error
}

type service struct {
	repo     Repository
	eventBus event.Bus
}

// NewService creates a new product service.
func NewService(repo Repository, bus event.Bus) Service {
	return &service{
		repo:     repo,
		eventBus: bus,
	}
}

// CreateProduct handles the business logic of creating a new product.
func (s *service) CreateProduct(ctx context.Context, orgID, name, description string, spec map[string]string) (*Product, error) {
	p, err := NewProduct(orgID, name, description, spec)
	if err != nil {
		return nil, err
	}
	return p, s.repo.Create(ctx, p)
}

// GetProductByID retrieves a product by its ID for a given organization.
func (s *service) GetProductByID(ctx context.Context, orgID, productID string) (*Product, error) {
	return s.repo.FindByID(ctx, orgID, productID)
}

// ListProductsForOrg lists all products for a given organization.
func (s *service) ListProductsForOrg(ctx context.Context, orgID string) ([]*Product, error) {
	return s.repo.ListByOrg(ctx, orgID)
}

// AppendLifecycle appends a new lifecycle event to a product instance.
func (s *service) AppendLifecycle(ctx context.Context, event LifecycleEvent) error {
	// Here you might add validation logic before saving.
	return s.repo.SaveLifecycleEvent(ctx, event)
}