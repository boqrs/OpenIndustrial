package workorder

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// Response defines the standard structure for a work order API response.
type Response struct {
	ID               uint                    `json:"id"`
	TenantID         string                  `json:"tenant_id"`
	FactoryID        uint                    `json:"factory_id"`
	ProductionLineID uint                    `json:"production_line_id"`
	ProductionPlanID uint                    `json:"production_plan_id"`
	ProductID        uint                    `json:"product_id"`
	BOMID            uint                    `json:"bom_id"`
	RoutingID        uint                    `json:"routing_id"`
	Code             string                  `json:"code"`
	PlannedQuantity  int64                   `json:"planned_quantity"`
	Priority         int                     `json:"priority"`
	DueDate          *time.Time              `json:"due_date,omitempty"`
	Status           model.WorkOrderStatus   `json:"status"`
	StartedAt        *time.Time              `json:"started_at,omitempty"`
	CompletedAt      *time.Time              `json:"completed_at,omitempty"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

// ToResponse converts a model.WorkOrder entity to a Response DTO.
func ToResponse(wo *model.WorkOrder) *Response {
	if wo == nil {
		return nil
	}
	return &Response{
		ID:               wo.ID,
		TenantID:         wo.TenantID.String(),
		FactoryID:        wo.FactoryID,
		ProductionLineID: wo.ProductionLineID,
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

// ToListResponse converts a slice of model.WorkOrder entities to a slice of Response DTOs.
func ToListResponse(wos []*model.WorkOrder) []*Response {
	if wos == nil {
		return nil
	}
	res := make([]*Response, 0, len(wos))
	for _, wo := range wos {
		res = append(res, ToResponse(wo))
	}
	return res
}