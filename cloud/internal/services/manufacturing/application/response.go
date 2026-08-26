package application

import (
	//"OpenIndustrial/cloud/internal/manufacturing/execution"
	//"OpenIndustrial/cloud/internal/manufacturing/routing"
	//"OpenIndustrial/cloud/internal/manufacturing/workorder"
	//"OpenIndustrial/cloud/internal/persistence"
	//"OpenIndustrial/cloud/internal/pkg/model"
	//"context"
	//"errors"
	//"fmt"
	"time"

	"github.com/google/uuid"
	//"gorm.io/gorm"
)

type CreateProductionExecutionResponse struct {
	ExecutionID uuid.UUID           `json:"execution_id"`
	WorkOrderID uuid.UUID           `json:"work_order_id"`
	ProductID   uuid.UUID           `json:"product_id"`
	RoutingID   uuid.UUID           `json:"routing_id"`
	Status      string              `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	Operations  []OperationResponse `json:"operations"`
}

type OperationResponse struct {
	OperationID uuid.UUID `json:"operation_id"`
	Sequence    int       `json:"sequence"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
}
