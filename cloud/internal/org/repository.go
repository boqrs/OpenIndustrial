package org

import "context"

// Repository defines the interface for interacting with organization persistence.
// This abstraction allows the service layer to be independent of the database implementation.
type Repository interface {
	// Create saves a new organization to the persistence layer.
	Create(ctx context.Context, org *Organization) error

	// FindByID retrieves an organization by its unique ID.
	FindByID(ctx context.Context, id string) (*Organization, error)

	// Update modifies an existing organization.
	Update(ctx context.Context, org *Organization) error

	// Delete removes an organization by its ID.
	Delete(ctx context.Context, id string) error

	// FindByParentID retrieves all organizations with a given parent ID.
	FindByParentID(ctx context.Context, parentID string) ([]*Organization, error)
}