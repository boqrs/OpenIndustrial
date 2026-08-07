package device

import (
	"time"

	"github.com/google/uuid"
)

// Device represents a physical device connected to a gateway.
type Device struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	GatewayID string    `json:"gateway_id"` // The gateway this device is connected through
	Name      string    `json:"name"`
	Model     string    `json:"model"`
	Status    string    `json:"status"` // e.g., "online", "offline", "error"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewDevice creates a new Device entity.
func NewDevice(orgID, gatewayID, name, model string) (*Device, error) {
	now := time.Now().UTC()
	return &Device{
		ID:        uuid.NewString(),
		OrgID:     orgID,
		GatewayID: gatewayID,
		Name:      name,
		Model:     model,
		Status:    "offline", // Initial status
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}