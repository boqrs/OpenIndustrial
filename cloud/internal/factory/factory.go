package factory

import (
	"context"

	"github.com/OpenIndustrial/cloud/internal/persistence/model"
)

// Repository defines the database operations for the Factory model.
type Repository interface {
	// Create inserts a new factory record into the database.
	Create(ctx context.Context, factory *model.Factory) error
}