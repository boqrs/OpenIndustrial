package workorder

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

type Response struct {
	ID               uint                    `json:"id"`
	ResourceUUID     uuid.UUID               `json:"resource_uuid"`
	TenantID         uuid.UUID               `json:"tenant_id"`
	ProductionPlanID uint                    `json:"production_plan_id"`
	ProductID        uint                    `json:"product_id"`
	BOMID            uint                    `json:"bom_id"`
	RoutingID        uint                    `json:"routing_id"`
	Code             string                  `json:"code"`
	PlannedQuantity  int64                   `json:"planned_quantity"`
	Priority         int                     `json:"priority"`
	DueDate          *time.Time              `json:"due_date"`
	Status           model.WorkOrderStatus   `json:"status"`
	StartedAt        *time.Time              `json:"started_at,omitempty"`
	CompletedAt      *time.Time              `json:"completed_at,omitempty"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

func ToResponse(wo *model.WorkOrder) *Response {
	if wo == nil {
		return nil
	}
	return &Response{
		ID:               wo.ID,
		ResourceUUID:     wo.ResourceUUID,
		TenantID:         wo.TenantID,
		ProductionPlanID: wo.ProductionPlanID,
		ProductID:        wo.ProductID,
		BOMID:            wo.BOMID,
		RoutingID:        wo.RoutingID,
		Code:             wo.Code,
		PlannedQuantity:  wo.PlannedQuantity,
		Priority:         wo.Priority,
		DueDate:          wo.DueDate,
		Status:           wo.Status,
		StartedAt:        wo.StartedAt,
		CompletedAt:      wo.CompletedAt,
		CreatedAt:        wo.CreatedAt,
		UpdatedAt:        wo.UpdatedAt,
	}
}

func ToResponses(wos []*model.WorkOrder) []*Response {
	responses := make([]*Response, 0, len(wos))
	for _, wo := range wos {
		responses = append(responses, ToResponse(wo))
	}
	return responses
}