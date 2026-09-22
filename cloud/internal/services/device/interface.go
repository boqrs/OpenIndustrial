package device

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

// Repository defines the persistence interface for devices.
type Repository interface {
	Create(
		ctx context.Context,
		entity *model.Device,
	) error

	CreateTx(
		ctx context.Context,
		entity *model.Device,
	) error

	GetByID(
		ctx context.Context,
		id uint,
	) (*model.Device, error)

	GetByResourceID(
		ctx context.Context,
		resourceID uint,
	) (*model.Device, error)

	GetBySerialNumber(
		ctx context.Context,
		serialNumber string,
	) (*model.Device, error)

	ActivateBySerialNumber(
		ctx context.Context,
		serialNumber string,
		customerUUID uuid.UUID,
	) (*model.Device, error)

	List(
		ctx context.Context,
		req *ListDevicesRequest,
	) ([]*model.Device, int64, error)

	Update(
		ctx context.Context,
		entity *model.Device,
	) error
}

type Service interface {
	GetDevice(
		ctx context.Context,
		deviceID uint,
	) (*DeviceResponse, error)

	ListDevices(
		ctx context.Context,
		req *ListDevicesRequest,
	) (*ListDevicesResponse, error)

	UpdateDevice(
		ctx context.Context,
		deviceID uint,
		req *UpdateDeviceRequest,
	) (*DeviceResponse, error)

	CreateFromExecutionResultTx(
		ctx context.Context,
		req *CreateDeviceFromExecutionResultRequest,
	) (*DeviceResponse, error)

	ActivateDevice(
		ctx context.Context,
		req *ActivateDeviceRequest,
	) (*DeviceResponse, error)
}
