package gateway

import (
	"time"

	"github.com/google/uuid"
)

// Gateway represents an edge gateway device.
type Gateway struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Model       string    `json:"model"`
	Status      string    `json:"status"`
	// The ID of the resource this gateway is associated with, e.g., a factory or a building.
	ResourceID  string    `json:"resource_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastHeartbeat time.Time `json:"last_heartbeat,omitempty"`
}

// NewGateway creates a new Gateway instance.
func NewGateway(name, description, model, resourceID string) *Gateway {
	return &Gateway{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Model:       model,
		ResourceID:  resourceID,
		Status:      "Offline", // Default status
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}