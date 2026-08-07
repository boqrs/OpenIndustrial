package product

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// MemoryRepository is an in-memory implementation of the product repository.
type MemoryRepository struct {
	mu              sync.RWMutex
	products        map[string]*Product
	instancesBySN   map[string]*ProductInstance // Keyed by SerialNumber
	instancesByID   map[string]*ProductInstance // Keyed by ID
}

// NewMemoryRepository creates a new in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		products:        make(map[string]*Product),
		instancesBySN:   make(map[string]*ProductInstance),
		instancesByID:   make(map[string]*ProductInstance),
	}
}

// --- Product Methods ---

func (r *MemoryRepository) CreateProduct(ctx context.Context, product *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.products[product.ID.String()]; exists {
		return fmt.Errorf("product with id %s already exists", product.ID)
	}
	r.products[product.ID.String()] = product
	return nil
}

func (r *MemoryRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	prod, ok := r.products[id.String()]
	if !ok {
		return nil, fmt.Errorf("product with id %s not found", id)
	}
	return prod, nil
}

func (r *MemoryRepository) ListProductsByOrg(ctx context.Context, orgID uuid.UUID) ([]*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var orgProducts []*Product
	for _, p := range r.products {
		if p.OrgID == orgID {
			orgProducts = append(orgProducts, p)
		}
	}
	return orgProducts, nil
}

// --- ProductInstance Methods ---

func (r *MemoryRepository) CreateInstance(ctx context.Context, instance *ProductInstance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.instancesBySN[instance.SerialNumber]; exists {
		return fmt.Errorf("product instance with sn %s already exists", instance.SerialNumber)
	}
    if _, exists := r.instancesByID[instance.ID.String()]; exists {
		return fmt.Errorf("product instance with id %s already exists", instance.ID)
	}
	r.instancesBySN[instance.SerialNumber] = instance
	r.instancesByID[instance.ID.String()] = instance
	return nil
}

func (r *MemoryRepository) GetInstanceBySN(ctx context.Context, sn string) (*ProductInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	instance, ok := r.instancesBySN[sn]
	if !ok {
		return nil, fmt.Errorf("product instance with sn %s not found", sn)
	}
	return instance, nil
}

func (r *MemoryRepository) GetInstanceByID(ctx context.Context, id uuid.UUID) (*ProductInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	instance, ok := r.instancesByID[id.String()]
	if !ok {
		return nil, fmt.Errorf("product instance with id %s not found", id)
	}
	return instance, nil
}

// --- LifecycleEvent Methods ---

func (r *MemoryRepository) AppendLifecycleEvent(ctx context.Context, event *LifecycleEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	instance, ok := r.instancesByID[event.ProductInstanceID.String()]
	if !ok {
		return fmt.Errorf("product instance with id %s not found for lifecycle event", event.ProductInstanceID)
	}

	instance.LifecycleEvents = append(instance.LifecycleEvents, event)
	return nil
}