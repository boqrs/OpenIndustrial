package identity

import (
	"context"
	"fmt"
	"sync"
)

// MemoryRepository is an in-memory implementation of the Repository interface.
// It is used for testing and development purposes.
type MemoryRepository struct {
	users       map[string]*User
	memberships map[string]*Membership
	mu          sync.RWMutex
}

// NewMemoryRepository creates a new MemoryRepository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users:       make(map[string]*User),
		memberships: make(map[string]*Membership),
	}
}

// CreateUser saves a new user to the in-memory store.
func (r *MemoryRepository) CreateUser(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user with id %s already exists", user.ID)
	}
	r.users[user.ID] = user
	return nil
}

// GetUser retrieves a user by their ID from the in-memory store.
func (r *MemoryRepository) GetUser(ctx context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with id %s not found", id)
	}
	return user, nil
}

// CreateMembership saves a new membership to the in-memory store.
func (r *MemoryRepository) CreateMembership(ctx context.Context, member *Membership) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.memberships[member.ID]; exists {
		return fmt.Errorf("membership with id %s already exists", member.ID)
	}
	r.memberships[member.ID] = member
	return nil
}

// ListMemberships retrieves all memberships for a given user ID.
func (r *MemoryRepository) ListMemberships(ctx context.Context, userID string) ([]*Membership, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userMemberships []*Membership
	for _, member := range r.memberships {
		if member.UserID == userID {
			userMemberships = append(userMemberships, member)
		}
	}
	return userMemberships, nil
}