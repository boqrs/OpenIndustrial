package device

import (
	"time"
)

// CreateDeviceRequest defines the structure for a request to create a new device.
type CreateDeviceRequest struct {
	Name      string `json:"name" binding:"required"`
	GatewayID string `json:"gateway_id" binding:"required"`
	Model     string `json:"model"`
}

// DeviceResponse defines the structure for a response containing device details.
type DeviceResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OrgID     string    `json:"org_id"`
	GatewayID string    `json:"gateway_id"`
	Model     string    `json:"model"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ToDeviceResponse converts a Device entity to a DeviceResponse DTO.
func ToDeviceResponse(dev *Device) *DeviceResponse {
	return &DeviceResponse{
		ID:        dev.ID.String(),
		Name:      dev.Name,
		OrgID:     dev.OrgID.String(),
		GatewayID: dev.GatewayID.String(),
		Model:     dev.Model,
		Status:    dev.Status,
		CreatedAt: dev.CreatedAt,
	}
}

// ToDeviceListResponse converts a slice of Device entities to a slice of DeviceResponse DTOs.
func ToDeviceListResponse(devs []*Device) []*DeviceResponse {
	res := make([]*DeviceResponse, len(devs))
	for i, dev := range devs {
		res[i] = ToDeviceResponse(dev)
	}
	return res
}