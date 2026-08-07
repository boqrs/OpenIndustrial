package org

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for organization storage.
type Repository interface {
	Create(ctx context.Context, org *Organization) error
	Get(ctx context.Context, id uuid.UUID) (*Organization, error)
}