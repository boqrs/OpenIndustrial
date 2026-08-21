
package device

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

// Repository defines the persistence interface for devices.
type Repository interface {
	Create(ctx context.Context, entity *model.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Device, error)
	GetByResourceID(ctx context.Context, resourceID uuid.UUID) (*model.Device, error)
	GetBySerialNumber(ctx context.Context, serialNumber string) (*model.Device, error)
	List(ctx context.Context, req *ListDevicesRequest) ([]*model.Device, int64, error)
	Update(ctx context.Context, entity *model.Device) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service interface {
	CreateDevice(ctx context.Context, req *CreateDeviceRequest) (*BootstrapCredentialResponse, error)
	GetDevice(ctx context.Context, deviceID uuid.UUID) (*DeviceResponse, error)
	ListDevices(ctx context.Context, req *ListDevicesRequest) (*ListDevicesResponse, error)
	UpdateDevice(ctx context.Context, deviceID uuid.UUID, req *UpdateDeviceRequest) (*DeviceResponse, error)
	DeleteDevice(ctx context.Context, deviceID uuid.UUID) error
}