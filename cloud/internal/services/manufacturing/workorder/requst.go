package workorder

import (
	"time"

	"github.com/google/uuid"
)

type CreateWorkOrderRequest struct {
	OrderNo string `json:"order_no" binding:"required"`

	ProductionPlanID uuid.UUID `json:"production_plan_id" binding:"required"`

	ProductID uuid.UUID `json:"product_id" binding:"required"`

	FactoryID uuid.UUID `json:"factory_id" binding:"required"`

	PlannedQuantity int `json:"planned_quantity" binding:"required,gt=0"`

	PlannedStartAt time.Time `json:"planned_start_at" binding:"required"`

	PlannedEndAt time.Time `json:"planned_end_at" binding:"required"`

	Priority int `json:"priority"`

	Description string `json:"description"`
}

type UpdateWorkOrderRequest struct {
	PlannedQuantity *int `json:"planned_quantity"`

	PlannedStartAt *time.Time `json:"planned_start_at"`

	PlannedEndAt *time.Time `json:"planned_end_at"`

	Priority *int `json:"priority"`

	Description *string `json:"description"`
}