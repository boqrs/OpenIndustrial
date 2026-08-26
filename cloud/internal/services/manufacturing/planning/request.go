package planning

import (
	"time"

	//"github.com/google/uuid"
)

// CreateProductionPlanRequest defines the request body for creating a new production plan.
type CreateProductionPlanRequest struct {
	PlanNo          string    `json:"planNo"`
	ProductID       uint      `json:"productID"`
	FactoryID       uint      `json:"factoryID"`
	PlannedQuantity int64     `json:"plannedQuantity"`
	PlannedStartAt  time.Time `json:"plannedStartAt"`
	PlannedEndAt    time.Time `json:"plannedEndAt"`
	Description     string    `json:"description"`
}

// UpdateProductionPlanRequest defines the request body for updating a production plan.
type UpdateProductionPlanRequest struct {
	PlannedQuantity *int64     `json:"plannedQuantity,omitempty"`
	PlannedStartAt  *time.Time `json:"plannedStartAt,omitempty"`
	PlannedEndAt    *time.Time `json:"plannedEndAt,omitempty"`
	Description     *string    `json:"description,omitempty"`
}
