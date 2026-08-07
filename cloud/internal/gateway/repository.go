package gateway

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for gateway storage.
type Repository interface {
	Create(ctx context.Context, gw *Gateway) (*Gateway, error)
	Get(ctx context.Context, id uuid.UUID) (*Gateway, error)
	List(ctx context.Context) ([]*Gateway, error)
	Update(ctx context.Context, gw *Gateway) error
}