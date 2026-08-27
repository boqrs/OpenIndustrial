package workorder

import (
	"time"

	"github.com/google/uuid"
)

type CreateRequest struct {
	ProductionPlanID uint      `json:"production_plan_id" binding:"required"`
	ProductID        uint      `json:"product_id" binding:"required"`
	BOMID            uint      `json:"bom_id" binding:"required"`
	RoutingID        uint      `json:"routing_id" binding:"required"`
	Code             string    `json:"code" binding:"required"`
	PlannedQuantity  int64     `json:"planned_quantity" binding:"required"`
	Priority         int       `json:"priority"`
	DueDate          *time.Time `json:"due_date"`
}

type UpdateRequest struct {
	Code            string     `json:"code"`
	PlannedQuantity int64      `json:"planned_quantity"`
	Priority        int        `json:"priority"`
	DueDate         *time.Time `json:"due_date"`
}

type ListRequest struct {
	TenantID  uuid.UUID `json:"-"`
	ProductID uint `json:"product_id"`
	Page      int       `json:"page"`
	PageSize  int       `json:"page_size"`
}