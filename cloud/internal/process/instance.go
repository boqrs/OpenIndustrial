package process

import "time"

// ProcessInstanceStatus represents the lifecycle state of a process instance.
type ProcessInstanceStatus string

const (
	StatusCreated   ProcessInstanceStatus = "created"
	StatusRunning   ProcessInstanceStatus = "running"
	StatusCompleted ProcessInstanceStatus = "completed"
	StatusFailed    ProcessInstanceStatus = "failed"
	StatusCancelled ProcessInstanceStatus = "cancelled"
	StatusSuspended ProcessInstanceStatus = "suspended"
)

// ProcessInstance is a live, running instance of a ProcessDefinition,
// typically tied to a specific product or work order.
type ProcessInstance struct {
	ID                  string                `json:"id"`
	// ProductID links this instance to the specific resource being produced (e.g., a product with a serial number).
	ProductID           string                `json:"product_id"`
	ProcessDefinitionID string                `json:"process_definition_id"`
	CurrentNodeID       string                `json:"current_node_id,omitempty"`
	Status              ProcessInstanceStatus `json:"status"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}