package gateway

import (
	"time"
)

// RegisterRequest defines the structure for a gateway registration request.
type RegisterRequest struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// GatewayResponse defines the structure for a response containing gateway details.
type GatewayResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	OrgID      string    `json:"org_id"`
	Status     string    `json:"status"`
	LastSeenAt time.Time `json:"last_seen_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToGatewayResponse converts a Gateway entity to a GatewayResponse DTO.
func ToGatewayResponse(gw *Gateway) *GatewayResponse {
	return &GatewayResponse{
		ID:         gw.ID.String(),
		Name:       gw.Name,
		OrgID:      gw.OrgID.String(),
		Status:     gw.Status,
		LastSeenAt: gw.LastSeenAt,
		CreatedAt:  gw.CreatedAt,
	}
}

// ToGatewayListResponse converts a slice of Gateway entities to a slice of GatewayResponse DTOs.
func ToGatewayListResponse(gws []*Gateway) []*GatewayResponse {
	res := make([]*GatewayResponse, len(gws))
	for i, gw := range gws {
		res[i] = ToGatewayResponse(gw)
	}
	return res
}