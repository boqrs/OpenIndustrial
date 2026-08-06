package user

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for interacting with user persistence.
type Repository interface {
	// Create saves a new user to the persistence layer.
	Create(ctx context.Context, user *User) error

	// FindByID retrieves a user by their unique ID.
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// FindByEmail retrieves a user by their email address.
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByUsername retrieves a user by their username.
	FindByUsername(ctx context.Context, username string) (*User, error)

	// Update modifies an existing user.
	Update(ctx context.Context, user *User) error

	// Delete removes a user by their ID.
	Delete(ctx context.Context, id uuid.UUID) error
}