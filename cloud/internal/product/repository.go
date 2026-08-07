package product

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrProductNotFound is returned when a product is not found.
	ErrProductNotFound = errors.New("product not found")
	// ErrLifecycleEventNotFound is returned when a lifecycle event is not found.
	ErrLifecycleEventNotFound = errors.New("lifecycle event not found")
)

// memoryRepository is an in-memory implementation of the Repository interface.
type memoryRepository struct {
	mu         sync.RWMutex
	products   map[string]*Product
	lifecycles map[string][]LifecycleEvent // Keyed by ProductInstanceID
}

// NewMemoryRepository creates a new in-memory product repository.
func NewMemoryRepository() Repository {
	return &memoryRepository{
		products:   make(map[string]*Product),
		lifecycles: make(map[string][]LifecycleEvent),
	}
}

// Create saves a new product to the repository.
func (r *memoryRepository) Create(ctx context.Context, product *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.products[product.ID] = product
	return nil
}

// FindByID retrieves a product by its ID.
func (r *memoryRepository) FindByID(ctx context.Context, orgID, productID string) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.products {
		if p.ID == productID && p.OrgID == orgID {
			return p, nil
		}
	}
	return nil, ErrProductNotFound
}

// ListByOrg lists all products for a given organization.
func (r *memoryRepository) ListByOrg(ctx context.Context, orgID string) ([]*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var results []*Product
	for _, p := range r.products {
		if p.OrgID == orgID {
			results = append(results, p)
		}
	}
	return results, nil
}

// SaveLifecycleEvent saves a new lifecycle event for a product instance.
func (r *memoryRepository) SaveLifecycleEvent(ctx context.Context, event LifecycleEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lifecycles[event.ProductInstanceID] = append(r.lifecycles[event.ProductInstanceID], event)
	return nil
}