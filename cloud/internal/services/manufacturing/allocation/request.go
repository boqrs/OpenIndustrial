package allocation

type CreateRequest struct {
	ProductionPlanID  uint  `json:"productionPlanID"`
	SalesOrderItemID  uint  `json:"salesOrderItemID"`
	AllocatedQuantity int64 `json:"allocatedQuantity"`
}
