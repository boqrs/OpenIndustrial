package gateway

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for data access operations for a Gateway.
type Repository interface {
	Create(ctx context.Context, gw *Gateway) error
	FindByID(ctx context.Context, orgID, gatewayID uuid.UUID) (*Gateway, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]*Gateway, error)
	Update(ctx context.Context, gw *Gateway) error
	Delete(ctx context.Context, orgID, gatewayID uuid.UUID) error
}