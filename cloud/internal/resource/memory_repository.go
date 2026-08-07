package resource

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MemoryRepository is an in-memory implementation of the resource repository.
type MemoryRepository struct {
	mu        sync.RWMutex
	resources map[string]*Resource
	relations map[string]*Relation
}

// NewMemoryRepository creates a new in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		resources: make(map[string]*Resource),
		relations: make(map[string]*Relation),
	}
}

// CreateResource saves a new resource.
func (r *MemoryRepository) CreateResource(ctx context.Context, resource *Resource) (*Resource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	resource.ID = uuid.New().String()
	resource.CreatedAt = time.Now().UTC()
	resource.UpdatedAt = resource.CreatedAt
	r.resources[resource.ID] = resource
	return resource, nil
}

// GetResource retrieves a resource by its ID.
func (r *MemoryRepository) GetResource(ctx context.Context, id string) (*Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resource, ok := r.resources[id]
	if !ok {
		return nil, fmt.Errorf("resource with id %s not found", id)
	}
	return resource, nil
}

// CreateRelation saves a new relation.
func (r *MemoryRepository) CreateRelation(ctx context.Context, relation *Relation) (*Relation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	relation.ID = uuid.New().String()
	relation.CreatedAt = time.Now().UTC()
	r.relations[relation.ID] = relation
	return relation, nil
}

// ListRelations retrieves relations based on query parameters.
// This is a simplified implementation. A real one would support more complex queries.
func (r *MemoryRepository) ListRelations(ctx context.Context, fromID, toID, relType string) ([]*Relation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*Relation
	for _, relation := range r.relations {
		match := true
		if fromID != "" && relation.FromID != fromID {
			match = false
		}
		if toID != "" && relation.ToID != toID {
			match = false
		}
		// Note: This checks for exact match of RelationType's string representation.
		if relType != "" && string(relation.Type) != relType {
			match = false
		}
		if match {
			results = append(results, relation)
		}
	}
	return results, nil
}

// UpdateResource updates an existing resource.
func (r *MemoryRepository) UpdateResource(ctx context.Context, resource *Resource) (*Resource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.resources[resource.ID]
	if !ok {
		return nil, fmt.Errorf("resource with id %s not found", resource.ID)
	}
	resource.UpdatedAt = time.Now().UTC()
	r.resources[resource.ID] = resource
	return resource, nil
}

// DeleteResource deletes a resource by its ID.
func (r *MemoryRepository) DeleteResource(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.resources[id]; !ok {
		return fmt.Errorf("resource with id %s not found", id)
	}
	delete(r.resources, id)
	// Also delete related relations
	for relID, rel := range r.relations {
		if rel.FromID == id || rel.ToID == id {
			delete(r.relations, relID)
		}
	}
	return nil
}