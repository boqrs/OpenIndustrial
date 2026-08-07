package permission

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// MemoryRepository is an in-memory implementation of the permission repository.
type MemoryRepository struct {
	mu          sync.RWMutex
	permissions map[string]*Permission
	policies    []*Policy // Policies are just tuples, store them in a slice
}

// NewMemoryRepository creates a new in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		permissions: make(map[string]*Permission),
		policies:    make([]*Policy, 0),
	}
}

// --- Permission Methods ---

func (r *MemoryRepository) GetPermission(ctx context.Context, id uuid.UUID) (*Permission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.permissions[id.String()]
	if !ok {
		return nil, fmt.Errorf("permission %s not found", id)
	}
	return p, nil
}

func (r *MemoryRepository) ListAllPermissions(ctx context.Context) ([]*Permission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*Permission, 0, len(r.permissions))
	for _, p := range r.permissions {
		list = append(list, p)
	}
	return list, nil
}

// --- Policy Methods ---

func (r *MemoryRepository) AddPolicy(ctx context.Context, policy *Policy) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Check for duplicates before adding
	for _, p := range r.policies {
		if p.Subject == policy.Subject && p.Object == policy.Object && p.Action == policy.Action {
			return nil // Already exists
		}
	}
	r.policies = append(r.policies, policy)
	return nil
}

func (r *MemoryRepository) RemovePolicy(ctx context.Context, policy *Policy) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, p := range r.policies {
		if p.Subject == policy.Subject && p.Object == policy.Object && p.Action == policy.Action {
			r.policies = append(r.policies[:i], r.policies[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("policy not found")
}

func (r *MemoryRepository) GetPoliciesForSubject(ctx context.Context, subject string) ([]*Policy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	results := make([]*Policy, 0)
	for _, p := range r.policies {
		if p.Subject == subject {
			results = append(results, p)
		}
	}
	return results, nil
}