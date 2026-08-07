package process

import "time"

// TaskStatus represents the lifecycle state of a task.
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusSkipped    TaskStatus = "skipped"
)

// Task represents the execution of a single ProcessNode within a ProcessInstance.
// This is the unit of work that is assigned and tracked.
type Task struct {
	ID                string     `json:"id"`
	ProcessInstanceID string     `json:"process_instance_id"`
	NodeID            string     `json:"node_id"`
	// ExecutorID could be a UserID for manual tasks or a SystemID for automatic ones.
	ExecutorID        string     `json:"executor_id,omitempty"`
	Status            TaskStatus `json:"status"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	// ResultData stores any output or results from the task execution.
	ResultData        map[string]interface{} `json:"result_data,omitempty"`
}