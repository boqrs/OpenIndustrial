package identity

import "context"

// Repository defines the interface for storing and retrieving identity-related entities.
type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id string) (*User, error)
	CreateMembership(ctx context.Context, member *Membership) error
	ListMemberships(ctx context.Context, userID string) ([]*Membership, error)
}