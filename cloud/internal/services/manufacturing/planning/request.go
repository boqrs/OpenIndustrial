package planning

import (
	"time"

	"github.com/google/uuid"
)

// CreateProductionPlanRequest defines the data needed to create a new production plan.
type CreateProductionPlanRequest struct {
	PlanNo          string    `json:"plan_no" binding:"required"`
	ProductID       uuid.UUID `json:"product_id" binding:"required"`
	FactoryID       uuid.UUID `json:"factory_id" binding:"required"`
	PlannedQuantity int       `json:"planned_quantity" binding:"required"`
	PlannedStartAt  time.Time `json:"planned_start_at" binding:"required"`
	PlannedEndAt    time.Time `json:"planned_end_at" binding:"required"`
	Description     string    `json:"description"`
}

// UpdateProductionPlanRequest defines the data for updating a production plan.
type UpdateProductionPlanRequest struct {
	PlannedQuantity *int       `json:"planned_quantity"`
	PlannedStartAt  *time.Time `json:"planned_start_at"`
	PlannedEndAt    *time.Time `json:"planned_end_at"`
	Description     *string    `json:"description"`
}
