package org

import "context"

// Repository defines the interface for data access operations for an Org.
// This interface is implemented by the persistence layer (e.g., a PostgreSQL implementation).
type Repository interface {
	// Create saves a new organization to the data store.
	Create(ctx context.Context, org *Org) error

	// FindByID retrieves an organization by its unique ID.
	FindByID(ctx context.Context, id string) (*Org, error)

	// Update modifies an existing organization in the data store.
	Update(ctx context.Context, org *Org) error

	// Delete removes an organization from the data store.
	Delete(ctx context.Context, id string) error
}