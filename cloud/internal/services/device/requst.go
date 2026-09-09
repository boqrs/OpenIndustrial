package device

import (
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
)

// CreateDeviceFromExecutionResultRequest defines the manufacturing
// information required to create a physical Device.
type CreateDeviceFromExecutionResultRequest struct {
	ProductID uint `json:"product_id" binding:"required"`

	WorkOrderID       uint `json:"work_order_id" binding:"required"`
	ExecutionID       uint `json:"execution_id" binding:"required"`
	ExecutionResultID uint `json:"execution_result_id" binding:"required"`

	SerialNumber string `json:"serial_number" binding:"required"`
	HardwareID   string `json:"hardware_id"`

	ParentResourceID *uint `json:"parent_resource_id"`
}

// UpdateDeviceRequest defines mutable device resource properties.
type UpdateDeviceRequest struct {
	Name             *string `json:"name"`
	ID               *uint   `json:"id"`
	ParentResourceID *uint   `json:"parent_resource_id"`
}

// ListDevicesRequest defines device filters and pagination.
type ListDevicesRequest struct {
	ProductID *uint               `json:"product_id,omitempty"`
	Status    *model.DeviceStatus `json:"status,omitempty"`
	ParentID  *uint               `json:"parent_id,omitempty"`
	pkg.BasePageReq
}
