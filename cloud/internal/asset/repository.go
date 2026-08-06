package asset

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for data access operations for an Asset.
type Repository interface {
	Create(ctx context.Context, asset *Asset) error
	FindByID(ctx context.Context, orgID, assetID uuid.UUID) (*Asset, error)
	FindBySN(ctx context.Context, orgID uuid.UUID, sn string) (*Asset, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]*Asset, error)
	Update(ctx context.Context, asset *Asset) error
}