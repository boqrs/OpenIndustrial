package application

import (
	"github.com/google/uuid"
)

type CreateProductionExecutionRequest struct {
	WorkOrderID uuid.UUID `json:"work_order_id"`
}

type operationResultEnvelope struct {
	Items []operationResultItem `json:"items"`
}

type operationResultItem struct {
	ItemKey string         `json:"item_key"`
	Data    map[string]any `json:"data"`
}

type productionItem struct {
	ItemKey string
	Data    map[string]any
}
