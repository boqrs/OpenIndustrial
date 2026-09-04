package executionresult


type CreateRequest struct {
	ExecutionID uint `json:"execution_id"`
	ProducedQuantity  int64 `json:"produced_quantity"`
	QualifiedQuantity int64 `json:"qualified_quantity"`
	RejectedQuantity  int64 `json:"rejected_quantity"`
}
