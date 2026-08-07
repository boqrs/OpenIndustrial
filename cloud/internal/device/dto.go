package device

import (
	"time"
)

// CreateDeviceRequest defines the request body for creating a new device.
type CreateDeviceRequest struct {
	Name string `json:"name" binding:"required"`
}

// DeviceResponse defines the standard response for a device.
type DeviceResponse struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToDeviceResponse converts a Device entity to a DeviceResponse.
func ToDeviceResponse(d *Device) *DeviceResponse {
	return &DeviceResponse{
		ID:        d.ID.String(),
		OrgID:     d.OrgID.String(),
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

// ToDeviceListResponse converts a slice of Device entities to a slice of DeviceResponse.
func ToDeviceListResponse(devices []*Device) []*DeviceResponse {
	res := make([]*DeviceResponse, len(devices))
	for i, d := range devices {
		res[i] = ToDeviceResponse(d)
	}
	return res
}