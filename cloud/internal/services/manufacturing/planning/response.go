package planning

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)
// ProductionPlanResponse is the structured response for a production plan.
type ProductionPlanResponse struct {
	ID              uuid.UUID                   `json:"id"`
	TenantID        uuid.UUID                   `json:"tenant_id"`
	PlanNo          string                      `json:"plan_no"`
	ProductID       uuid.UUID                   `json:"product_id"`
	FactoryID       uuid.UUID                   `json:"factory_id"`
	PlannedQuantity int                         `json:"planned_quantity"`
	PlannedStartAt  time.Time                   `json:"planned_start_at"`
	PlannedEndAt    time.Time                   `json:"planned_end_at"`
	Status          model.ProductionPlanStatus  `json:"status"`
	Description     string                      `json:"description"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
}