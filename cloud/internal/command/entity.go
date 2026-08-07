package command

import (
	"time"

	"github.com/google/uuid"
)

// CommandType defines the type of action the command represents.
type CommandType string

const (
	// CommandTypeRead requests to read data from a resource.
	CommandTypeRead CommandType = "Read"
	// CommandTypeWrite requests to write data or configuration to a resource.
	CommandTypeWrite CommandType = "Write"
	// CommandTypeExecute requests to execute a specific function or method on a resource.
	CommandTypeExecute CommandType = "Execute"
)

// CommandStatus tracks the lifecycle of a command.
type CommandStatus string

const (
	// StatusPending means the command is created but not yet sent.
	StatusPending CommandStatus = "Pending"
	// StatusSent means the command has been sent to the target resource.
	StatusSent CommandStatus = "Sent"
	// StatusAcknowledged means the target has acknowledged receipt of the command.
	StatusAcknowledged CommandStatus = "Acknowledged"
	// StatusSuccess means the command was executed successfully.
	StatusSuccess CommandStatus = "Success"
	// StatusFailed means the command execution failed.
	StatusFailed CommandStatus = "Failed"
	// StatusTimeout means the command timed out without a response.
	StatusTimeout CommandStatus = "Timeout"
)

// Command represents a request sent from the cloud to a resource to perform an action.
type Command struct {
	ID         uuid.UUID `json:"id"`
	ResourceID uuid.UUID `json:"resource_id"` // The target resource for the command
	Type       CommandType `json:"type"`
	// Parameters for the command, e.g., {"register": "0x1001", "value": 1}.
	Parameters map[string]interface{} `json:"parameters"`
	Status     CommandStatus `json:"status"`
	Timeout    time.Duration `json:"timeout"` // Duration in seconds
	IssuedAt   time.Time `json:"issued_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	// Response from the device, if any.
	Response map[string]interface{} `json:"response,omitempty"`
}