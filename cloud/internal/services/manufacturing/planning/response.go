package planning

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

// ProductionPlanResponse defines the response body for a production plan.
// ProductionPlanResponse defines the response body for a production plan.
type ProductionPlanResponse struct {
	ID              uint                       `json:"id"`
	ResourceID    uint                  `json:"resourceUuid"`
	TenantID        uuid.UUID                  `json:"tenantId"`
	PlanNo          string                     `json:"planNo"`
	ProductID       uint                       `json:"productId"`
	FactoryID       uint                       `json:"factoryId"`
	PlannedQuantity int64                      `json:"plannedQuantity"`
	PlannedStartAt  time.Time                  `json:"plannedStartAt"`
	PlannedEndAt    time.Time                  `json:"plannedEndAt"`
	Status          model.ProductionPlanStatus `json:"status"`
	Description     string                     `json:"description"`
	CreatedAt       time.Time                  `json:"createdAt"`
	UpdatedAt       time.Time                  `json:"updatedAt"`
}