package workorder

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/google/uuid"
)

// CreateRequest defines the structure for creating a new work order.
// FactoryID is intentionally omitted; it will be derived from the ProductionPlan.
type CreateRequest struct {
	ProductionPlanID uint       `json:"production_plan_id" binding:"required"`
	ProductionLineID uint       `json:"production_line_id" binding:"required"`
	ProductID        uint       `json:"product_id" binding:"required"`
	BOMID            uint       `json:"bom_id" binding:"required"`
	RoutingID        uint       `json:"routing_id" binding:"required"`
	Code             string     `json:"code" binding:"required"`
	PlannedQuantity  int64      `json:"planned_quantity" binding:"required"`
	Priority         int        `json:"priority"`
	DueDate          *time.Time `json:"due_date"`
}


// UpdateRequest defines the structure for updating an existing work order.
// Only certain fields of a work order in 'draft' status can be updated.
type UpdateRequest struct {
	Code            string     `json:"code" binding:"required"`
	PlannedQuantity int64      `json:"planned_quantity" binding:"required"`
	Priority        int        `json:"priority"`
	DueDate         *time.Time `json:"due_date"`
}
type ListRequest struct {
	TenantID  uuid.UUID `json:"-"`
	ProductID uint `json:"product_id"`
	pkg.BasePageReq
}