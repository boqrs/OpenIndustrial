package param

import (
	"github.com/google/uuid"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"

)

// CreateDeviceRequest defines the payload for creating a new device.
type CreateDeviceRequest struct {
	Name           string    `json:"name" binding:"required"`
	ProductModelID uuid.UUID `json:"product_model_id" binding:"required"`
	SerialNumber   string    `json:"serial_number"`
	HardwareID     string    `json:"hardware_id"`
	ParentResourceID *uuid.UUID `json:"parent_resource_id"` // For placing the device in the resource tree
}

// UpdateDeviceRequest defines the payload for updating an existing device.
type UpdateDeviceRequest struct {
	Name             *string    `json:"name"`
	ParentResourceID *uuid.UUID `json:"parent_resource_id"`
}

// ListDevicesRequest defines the filters and pagination for listing devices.
type ListDevicesRequest struct {
	Page           int
	PageSize       int
	ProductModelID *uuid.UUID
	Status         *model.DeviceStatus
	ParentID       *uuid.UUID
}