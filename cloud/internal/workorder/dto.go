package workorder

import "time"

// CreateWorkOrderRequest defines the request body for creating a new work order.
type CreateWorkOrderRequest struct {
	ProductID   string    `json:"product_id" binding:"required"`
	Quantity    int       `json:"quantity" binding:"required,gt=0"`
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

// ToWorkOrderResponse converts a WorkOrder entity to a WorkOrderResponse.
func ToWorkOrderResponse(wo *WorkOrder) *WorkOrderResponse {
	return &WorkOrderResponse{
		ID:          wo.ID,
		OrgID:       wo.OrgID,
		ProductID:   wo.ProductID,
		Quantity:    wo.Quantity,
		Status:      wo.Status,
		ScheduledAt: wo.ScheduledAt,
		StartedAt:   wo.StartedAt,
		CompletedAt: wo.CompletedAt,
		CreatedAt:   wo.CreatedAt,
		UpdatedAt:   wo.UpdatedAt,
	}
}

// ToWorkOrderListResponse converts a slice of WorkOrder entities to a slice of WorkOrderResponse.
func ToWorkOrderListResponse(wos []*WorkOrder) []*WorkOrderResponse {
	res := make([]*WorkOrderResponse, len(wos))
	for i, wo := range wos {
		res[i] = ToWorkOrderResponse(wo)
	}
	return res
}