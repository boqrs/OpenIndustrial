package execution

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

type ExecutionResponse struct {
	ID           uint                            `json:"id"`
	ResourceUUID uuid.UUID                       `json:"resourceUuid"`
	TenantID     uuid.UUID                       `json:"tenantId"`
	WorkOrderID  uint                            `json:"workOrderId"`
	DeviceID     *uint                           `json:"deviceId,omitempty"`
	Quantity     int64                           `json:"quantity"`
	Status       model.ProductionExecutionStatus `json:"status"`
	StartedAt    *time.Time                      `json:"startedAt,omitempty"`
	CompletedAt  *time.Time                      `json:"completedAt,omitempty"`
	CreatedAt    time.Time                       `json:"createdAt"`
	UpdatedAt    time.Time                       `json:"updatedAt"`
}

// --- DTOs for Operation ---

type OperationResponse struct {
	ID             uint                       `json:"id"`
	ExecutionID    uint                            `json:"executionId"`
	Sequence       int                             `json:"sequence"`
	WorkCenterID   uint                            `json:"workCenterId"`
	Status         model.ExecutionOperationStatus  `json:"status"`
	Result         map[string]any                  `json:"result,omitempty"`
	StartedAt      *time.Time                      `json:"startedAt,omitempty"`
	CompletedAt    *time.Time                      `json:"completedAt,omitempty"`
	CreatedAt      time.Time                       `json:"createdAt"`
	UpdatedAt      time.Time                       `json:"updatedAt"`
}