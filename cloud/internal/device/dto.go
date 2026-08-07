package device

import "time"

// CreateDeviceRequest defines the request body for creating a new device.
type CreateDeviceRequest struct {
	GatewayID string `json:"gateway_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Model     string `json:"model" binding:"required"`
}

// DeviceResponse defines the standard response for a device.
type DeviceResponse struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	GatewayID string    `json:"gateway_id"`
	Name      string    `json:"name"`
	Model     string    `json:"model"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToDeviceResponse converts a Device entity to a DeviceResponse.
func ToDeviceResponse(d *Device) *DeviceResponse {
	return &DeviceResponse{
		ID:        d.ID,
		OrgID:     d.OrgID,
		GatewayID: d.GatewayID,
		Name:      d.Name,
		Model:     d.Model,
		Status:    string(d.Status),
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