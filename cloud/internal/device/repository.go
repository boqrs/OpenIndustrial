package device

import (
	"context"

	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/google/uuid"
)

// Repository defines the persistence interface for devices.
type Repository interface {
	Create(ctx context.Context, entity *model.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Device, error)
	GetByResourceID(ctx context.Context, resourceID uuid.UUID) (*model.Device, error)
	GetBySerialNumber(ctx context.Context, serialNumber string) (*model.Device, error)
	List(ctx context.Context, req *param.ListDevicesRequest) ([]*model.Device, int64, error)
	Update(ctx context.Context, entity *model.Device) error
	Delete(ctx context.Context, id uuid.UUID) error
}