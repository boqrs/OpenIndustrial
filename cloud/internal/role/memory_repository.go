package role

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// MemoryRepository is an in-memory implementation of the role Repository.
type MemoryRepository struct {
	roles    map[string]*Role
	bindings map[string]*Binding
	mu       sync.RWMutex
}

// NewMemoryRepository creates a new MemoryRepository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		roles:    make(map[string]*Role),
		bindings: make(map[string]*Binding),
	}
}

func (r *MemoryRepository) CreateRole(ctx context.Context, role *Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.roles[role.ID.String()]; exists {
		return fmt.Errorf("role with id %s already exists", role.ID)
	}
	r.roles[role.ID.String()] = role
	return nil
}

func (r *MemoryRepository) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	role, ok := r.roles[roleID.String()]
	if !ok {
		return nil, fmt.Errorf("role with id %s not found", roleID)
	}
	return role, nil
}

func (r *MemoryRepository) ListRolesByOrg(ctx context.Context, orgID uuid.UUID) ([]*Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var orgRoles []*Role
	for _, role := range r.roles {
		if role.OrgID == orgID {
			orgRoles = append(orgRoles, role)
		}
	}
	return orgRoles, nil
}

func (r *MemoryRepository) BindRole(ctx context.Context, binding *Binding) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.bindings[binding.ID]; exists {
		return fmt.Errorf("binding with id %s already exists", binding.ID)
	}
	r.bindings[binding.ID] = binding
	return nil
}

func (r *MemoryRepository) ListRolesByMembership(ctx context.Context, membershipID string) ([]*Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var memberRoles []*Role
	for _, binding := range r.bindings {
		if binding.MembershipID == membershipID {
			role, exists := r.roles[binding.RoleID]
			if exists {
				memberRoles = append(memberRoles, role)
			}
		}
	}
	return memberRoles, nil
}