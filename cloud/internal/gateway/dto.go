package gateway

import (
	"time"

	"github.com/google/uuid"
)

// CreateGatewayRequest defines the request body for creating a new gateway.
type CreateGatewayRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Model       string `json:"model"`
	// The ID of the resource this gateway is associated with, e.g., a factory or a building.
	ResourceID string `json:"resource_id"`
}

// RegisterRequest is used for the initial registration of a gateway.
// It might contain authentication details like a pre-shared key.
type RegisterRequest struct {
	Model string `json:"model"`
	// ... other registration fields
}

// GatewayResponse defines the standard gateway representation for API responses.
type GatewayResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	Model         string    `json:"model"`
	Status        string    `json:"status"`
	ResourceID    string    `json:"resource_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastHeartbeat time.Time `json:"last_heartbeat,omitempty"`
}

// ToGatewayResponse converts a Gateway entity to a GatewayResponse DTO.
func ToGatewayResponse(gw *Gateway) *GatewayResponse {
	return &GatewayResponse{
		ID:            gw.ID,
		Name:          gw.Name,
		Description:   gw.Description,
		Model:         gw.Model,
		Status:        gw.Status,
		ResourceID:    gw.ResourceID,
		CreatedAt:     gw.CreatedAt,
		UpdatedAt:     gw.UpdatedAt,
		LastHeartbeat: gw.LastHeartbeat,
	}
}

// ToGatewayListResponse converts a slice of Gateway entities to a slice of GatewayResponse DTOs.
func ToGatewayListResponse(gws []*Gateway) []*GatewayResponse {
	responses := make([]*GatewayResponse, len(gws))
	for i, gw := range gws {
		responses[i] = ToGatewayResponse(gw)
	}
	return responses
}