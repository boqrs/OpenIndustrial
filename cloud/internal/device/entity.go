package device

import (
	"time"

	"github.com/google/uuid"
)

// Device represents a physical or logical device that collects data.
// It is connected to the platform via a Gateway.
type Device struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OrgID     uuid.UUID `json:"org_id"`
	GatewayID uuid.UUID `json:"gateway_id"` // Which gateway this device is connected through
	Model     string    `json:"model"`      // e.g., "S7-1200", "UR5e"
	Status    string    `json:"status"`    // e.g., "Online", "Offline", "Error"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewDevice creates a new Device entity.
func NewDevice(orgID, gatewayID uuid.UUID, name, model string) (*Device, error) {
	if name == "" {
		return nil, ErrDeviceNameRequired
	}

	now := time.Now().UTC()
	return &Device{
		ID:        uuid.New(),
		Name:      name,
		OrgID:     orgID,
		GatewayID: gatewayID,
		Model:     model,
		Status:    "Offline",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}