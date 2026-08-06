package device

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for data access operations for a Device.
type Repository interface {
	Create(ctx context.Context, dev *Device) error
	FindByID(ctx context.Context, orgID, deviceID uuid.UUID) (*Device, error)
	ListByGateway(ctx context.Context, orgID, gatewayID uuid.UUID) ([]*Device, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]*Device, error)
	Update(ctx context.Context, dev *Device) error
	Delete(ctx context.Context, orgID, deviceID uuid.UUID) error
}