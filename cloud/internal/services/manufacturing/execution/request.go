package execution

import "github.com/google/uuid"

type CreateExecutionRequest struct {
	WorkOrderID uint `json:"workOrderId"`

	// ResourceID is created by the manufacturing application layer.
	// Execution itself does not create the Resource because Resource creation
	// is a cross-domain application concern.
	ResourceID uint `json:"-"`
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
