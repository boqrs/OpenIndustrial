package workorder

import (
	"time"

	//"github.com/google/uuid"
)


// CreateWorkOrderRequest defines the request body for creating a new work order.
type CreateWorkOrderRequest struct {
	Code             string     `json:"code"`
	ProductionPlanID uint       `json:"productionPlanId"`
	ProductID        uint       `json:"productId"`
	RoutingID        uint       `json:"routingId"`
	PlannedQuantity  int64      `json:"plannedQuantity"`
	DueDate          *time.Time `json:"dueDate,omitempty"`
	Priority         int        `json:"priority"`
	Description      string     `json:"description"`
}

// UpdateWorkOrderRequest defines the request body for updating a work order.
type UpdateWorkOrderRequest struct {
	PlannedQuantity *int64     `json:"plannedQuantity,omitempty"`
	DueDate         *time.Time `json:"dueDate,omitempty"`
	Priority        *int       `json:"priority,omitempty"`
	Description     *string    `json:"description,omitempty"`
}