package workorder

import (
	"time"

	"github.com/google/uuid"
)

// CreateWorkOrderRequest defines the request body for creating a work order.
type CreateWorkOrderRequest struct {
	ProductID   uuid.UUID `json:"product_id" binding:"required"`
	Quantity    int       `json:"quantity" binding:"required"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

// WorkOrderResponse defines the standard response for a work order.
type WorkOrderResponse struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	ProductID   string    `json:"product_id"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"`
	ScheduledAt time.Time `json:"scheduled_at"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToWorkOrderResponse converts a WorkOrder entity to a WorkOrderResponse DTO.
func ToWorkOrderResponse(wo *WorkOrder) *WorkOrderResponse {
	return &WorkOrderResponse{
		ID:          wo.ID.String(),
		OrgID:       wo.OrgID.String(),
		ProductID:   wo.ProductID.String(),
		Quantity:    wo.Quantity,
		Status:      wo.Status,
		ScheduledAt: wo.ScheduledAt,
		StartedAt:   wo.StartedAt,
		CompletedAt: wo.CompletedAt,
		CreatedAt:   wo.CreatedAt,
		UpdatedAt:   wo.UpdatedAt,
	}
}

// ToWorkOrderListResponse converts a slice of WorkOrder entities to a slice of WorkOrderResponse DTOs.
func ToWorkOrderListResponse(wos []*WorkOrder) []*WorkOrderResponse {
	res := make([]*WorkOrderResponse, len(wos))
	for i, wo := range wos {
		res[i] = ToWorkOrderResponse(wo)
	}
	return res
}