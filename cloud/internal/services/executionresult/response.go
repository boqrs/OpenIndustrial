package executionresult

import (
	"time"
	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// ExecutionResultResponse represents the response structure for ExecutionResult.
type Response struct { 
	ID uint `json:"id"` 
	TenantID uuid.UUID `json:"tenant_id"` 
	ExecutionID uint `json:"execution_id"` 
	WorkOrderID uint `json:"work_order_id"`
	ProducedQuantity int64 `json:"produced_quantity"`
	QualifiedQuantity int64 `json:"qualified_quantity"` 
	RejectedQuantity int64 `json:"rejected_quantity"` 
	Status model.ExecutionResultStatus `json:"status"` 
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"` 
} 

func toResponse(entity *model.ExecutionResult) *Response { 
		if entity == nil {
			 return nil 
		}

		return &Response{ 
			ID: entity.ID,
			TenantID: entity.TenantID,
			ExecutionID: entity.ExecutionID,
			WorkOrderID: entity.WorkOrderID, 
			ProducedQuantity: entity.ProducedQuantity,
			QualifiedQuantity: entity.QualifiedQuantity,
			RejectedQuantity: entity.RejectedQuantity, 
			Status: entity.Status, 
			ConfirmedAt: entity.ConfirmedAt,
			CreatedAt: entity.CreatedAt, 
			UpdatedAt: entity.UpdatedAt,
		}
}