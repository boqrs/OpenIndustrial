package device

import (
	"context"

	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

type DeviceTypeRepository interface {
	Create(ctx context.Context, entity *model.DeviceType) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.DeviceType, error)
	GetByResourceID(ctx context.Context, resourceID uuid.UUID) (*model.DeviceType, error)
	GetByCode(ctx context.Context, code string) (*model.DeviceType, error)
	List(ctx context.Context) ([]*model.DeviceType, error)
	Update(ctx context.Context, entity *model.DeviceType) error
}

type DeviceRepository interface {
	Create(ctx context.Context, entity *model.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Device, error)
	GetByResourceID(ctx context.Context, resourceID uuid.UUID) (*model.Device, error)
	List(ctx context.Context, deviceTypeID *uuid.UUID) ([]*model.Device, error)
	Update(ctx context.Context, entity *model.Device) error
	Delete(ctx context.Context, id uuid.UUID) error
}