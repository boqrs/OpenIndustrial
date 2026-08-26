package workorder

import (
	"time"

	"github.com/google/uuid"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// WorkOrderResponse defines the response body for a work order.
type WorkOrderResponse struct {
	ID               uint                  `json:"id"`
	ResourceUUID     uuid.UUID             `json:"resourceUuid"`
	TenantID         uuid.UUID             `json:"tenantId"`
	Code             string                `json:"code"`
	ProductionPlanID uint                  `json:"productionPlanId"`
	ProductID        uint                  `json:"productId"`
	RoutingID        uint                  `json:"routingId"`
	PlannedQuantity  int64                 `json:"plannedQuantity"`
	DueDate          *time.Time            `json:"dueDate,omitempty"`
	Status           model.WorkOrderStatus `json:"status"`
	Priority         int                   `json:"priority"`
	Description      string                `json:"description"`
	StartedAt        *time.Time            `json:"startedAt,omitempty"`
	CompletedAt      *time.Time            `json:"completedAt,omitempty"`
	CreatedAt        time.Time             `json:"createdAt"`
	UpdatedAt        time.Time             `json:"updatedAt"`
}