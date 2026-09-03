package execution

import "github.com/google/uuid"

type CreateExecutionRequest struct {
	WorkOrderID uint  `json:"workOrderId"`
	DeviceID    *uint `json:"deviceId,omitempty"`
}


type RoutingOperationSnapshot struct {
	ID                      uuid.UUID
	Sequence                int
	Code                    string
	Name                    string
	Description             string
	WorkstationID           uint
	StandardDurationSeconds int
	Required                bool
	Parameters              map[string]any
}