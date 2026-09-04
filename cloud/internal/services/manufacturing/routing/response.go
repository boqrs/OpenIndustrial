package routing

import (
	"time"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// --- Response Structs ---
type OperationResponse struct {
	ID             uint `json:"id"`
	RoutingID      uint      `json:"routingId"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	WorkStationID   uint      `json:"workstationId"`
	Sequence       int       `json:"sequence"`
	SetupTime      int       `json:"setupTime"`
	ProcessingTime int       `json:"processingTime"`
	Parameters     []byte `json:"parameters"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type RoutingResponse struct {
	ID           uint                `json:"id"`
	ResourceID uint           `json:"resourceUuid"`
	TenantID     uuid.UUID           `json:"tenantId"`
	ProductID    uint                `json:"productId"`
	Name         string              `json:"name"`
	Version      int              `json:"version"`
	Description  string              `json:"description"`
	Status       model.RoutingStatus `json:"status"`
	IsDefault    bool                `json:"isDefault"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
}
