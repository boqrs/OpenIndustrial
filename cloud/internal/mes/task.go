package mes

import "time"

// StationTask is an instance of a ProcessStep assigned to a specific ProductInstance and Station.
type StationTask struct {
	ID                string
	ProductInstanceID string // Foreign key to product.ProductInstance
	ProcessStepID     string // Foreign key to ProcessStep
	StationID         string // The station assigned to execute the task
	Status            TaskStatus
	OperatorID        string // The user who is performing the task
	CreatedAt         time.Time
	StartedAt         *time.Time
	CompletedAt       *time.Time
}

// TaskStatus defines the lifecycle of a station task.
type TaskStatus string

const (
	TaskCreated TaskStatus = "created"
	TaskRunning TaskStatus = "running"
	TaskPassed  TaskStatus = "passed"
	TaskFailed  TaskStatus = "failed"
	TaskBlocked TaskStatus = "blocked" // e.g., waiting for materials or a quality check
)