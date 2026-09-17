package allocation

import "time"

type Response struct {
	ID uint `json:"id"`

	ProductionPlanID uint `json:"productionPlanID"`
	SalesOrderItemID uint `json:"salesOrderItemID"`

	AllocatedQuantity int64 `json:"allocatedQuantity"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
