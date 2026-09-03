package execution

import "github.com/google/uuid"

type CreateExecutionRequest struct {
	WorkOrderID uint  `json:"workOrderId"`
	DeviceID    *uint `json:"deviceId,omitempty"`
	//Quantity    int64 `json:"quantity"`
}


type RoutingOperationSnapshot struct {
	ID                      uuid.UUID
	Sequence                int
	Code                    string
	Name                    string
	Description             string
	WorkstationID           *uuid.UUID
	StandardDurationSeconds int
	Required                bool
	Parameters              map[string]any
}