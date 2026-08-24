package workorder

import (
	"time"

	"github.com/google/uuid"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type WorkOrderResponse struct {
	ID uuid.UUID `json:"id"`

	TenantID uuid.UUID `json:"tenant_id"`

	OrderNo string `json:"order_no"`

	ProductionPlanID uuid.UUID `json:"production_plan_id"`

	ProductID uuid.UUID `json:"product_id"`

	FactoryID uuid.UUID `json:"factory_id"`

	PlannedQuantity int `json:"planned_quantity"`

	CompletedQuantity int `json:"completed_quantity"`

	RemainingQuantity int `json:"remaining_quantity"`

	PlannedStartAt time.Time `json:"planned_start_at"`

	PlannedEndAt time.Time `json:"planned_end_at"`

	Status model.WorkOrderStatus `json:"status"`

	Priority int `json:"priority"`

	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}