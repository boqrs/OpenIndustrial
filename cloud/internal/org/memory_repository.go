package org

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// MemoryRepository is an in-memory implementation of the org repository.
type MemoryRepository struct {
	mu   sync.RWMutex
	orgs map[string]*Organization
}

// NewMemoryRepository creates a new in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		orgs: make(map[string]*Organization),
	}
}

// Create saves a new organization.
func (r *MemoryRepository) Create(ctx context.Context, org *Organization) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orgs[org.ID.String()]; exists {
		return fmt.Errorf("organization with id %s already exists", org.ID)
	}
	r.orgs[org.ID.String()] = org
	return nil
}

// Get retrieves an organization by its ID.
func (r *MemoryRepository) Get(ctx context.Context, id uuid.UUID) (*Organization, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	org, ok := r.orgs[id.String()]
	if !ok {
		return nil, fmt.Errorf("organization with id %s not found", id)
	}
	return org, nil
}