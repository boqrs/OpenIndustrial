package device

import (
	"time"
	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)


// DeviceResponse is the standard representation of a device returned by the API.
type DeviceResponse struct {
	ID               uuid.UUID    `json:"id"`
	ResourceID       uuid.UUID    `json:"resource_id"`
	ProductID   uint   `json:"product_model_id"`
	Name             string       `json:"name"`
	SerialNumber     string       `json:"serial_number"`
	HardwareID       string       `json:"hardware_id"`
	Status           model.DeviceStatus `json:"status"`
	ParentResourceID *uuid.UUID   `json:"parent_resource_id"`
	CreatedAt        string       `json:"created_at"`
	UpdatedAt        string       `json:"updated_at"`
	LastOnlineAt     *string      `json:"last_online_at,omitempty"`
}

// ListDevicesResponse is the paginated response for a list of devices.
type ListDevicesResponse struct {
	Items      []*DeviceResponse `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
}

type BootstrapCredentialResponse struct {
	ResourceID   uuid.UUID `json:"resource_id"`
	CredentialID uuid.UUID `json:"credential_id"`

	// 明文 Token 只返回一次。
	// 数据库绝对不保存这个值。
	Token string `json:"token"`

	CreatedAt time.Time `json:"created_at"`
}

