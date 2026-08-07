package mes

import (
	"context"
)

// Repository defines the storage interface for MES tasks.
type Repository interface {
	// GetCurrentTaskForProduct retrieves the current in-progress or pending task for a specific product instance ID.
	GetCurrentTaskForProduct(ctx context.Context, productInstanceID string) (*Task, error)
}