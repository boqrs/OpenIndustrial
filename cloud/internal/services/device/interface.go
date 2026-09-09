
package device

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// Repository defines the persistence interface for devices.
type Repository interface {
    Create(ctx context.Context, entity *model.Device) error
	CreateTx(ctx context.Context, entity *model.Device) error
	CreateBatchTx(ctx context.Context,devices []*model.Device) error
    GetByID(ctx context.Context, id uint) (*model.Device, error)
    GetByResourceID(ctx context.Context, resourceID uint) (*model.Device, error)
    GetBySerialNumber(ctx context.Context, serialNumber string) (*model.Device, error)
    List( ctx context.Context,req *ListDevicesRequest) ([]*model.Device, int64, error)
    Update(ctx context.Context, entity *model.Device) error
    Delete(ctx context.Context, id uint) error
}

type Service interface {
    GetDevice(ctx context.Context,deviceID uint) (*DeviceResponse, error)

    ListDevices(ctx context.Context,req *ListDevicesRequest) (*ListDevicesResponse, error)

    UpdateDevice(ctx context.Context,deviceID uint,req *UpdateDeviceRequest) (*DeviceResponse, error)

    DeleteDevice(ctx context.Context,deviceID uint) error

    CreateFromExecutionResultTx(ctx context.Context,req *CreateDeviceFromExecutionResultRequest) (*DeviceResponse, error)
	CreateFromExecutionResultBatchTx(ctx context.Context,reqs []*CreateDeviceFromExecutionResultRequest) ([]*DeviceResponse, error)
}