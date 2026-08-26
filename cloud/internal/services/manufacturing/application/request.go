package application

import (
	"github.com/google/uuid"
)

type CreateProductionExecutionRequest struct {
	WorkOrderID uuid.UUID `json:"work_order_id"`
}
