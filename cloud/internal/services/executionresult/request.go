package executionresult

type CreateResultRequest struct {
	ExecutionID       uint  `json:"execution_id" binding:"required"`
	ProducedQuantity  int64 `json:"produced_quantity"`
	QualifiedQuantity int64 `json:"qualified_quantity"`
	RejectedQuantity  int64 `json:"rejected_quantity"`
	//WorkOrderID uint `json:"work_order_id" binding:"required"`
}
