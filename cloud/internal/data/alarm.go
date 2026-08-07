package data

import "time"

// AlarmLevel defines the severity of an alarm.
type AlarmLevel string

const (
	AlarmInfo     AlarmLevel = "INFO"
	AlarmWarning  AlarmLevel = "WARNING"
	AlarmError    AlarmLevel = "ERROR"
	AlarmCritical AlarmLevel = "CRITICAL"
)

// AlarmStatus tracks the lifecycle of an alarm.
type AlarmStatus string

const (
	StatusTriggered    AlarmStatus = "Triggered"
	StatusAcknowledged AlarmStatus = "Acknowledged"
	StatusResolved     AlarmStatus = "Resolved"
	StatusClosed       AlarmStatus = "Closed"
)

// Alarm represents a significant event triggered by a stream rule or other system logic.
type Alarm struct {
	ID         string      `json:"id"`
	ResourceID string      `json:"resource_id"`
	Level      AlarmLevel  `json:"level"`
	Message    string      `json:"message"`
	Status     AlarmStatus `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	ResolvedAt *time.Time  `json:"resolved_at,omitempty"`
}

// AlarmEvent is an immutable record of an action taken on an alarm,
// providing a full audit trail for alarm management.
type AlarmEvent struct {
	ID         string    `json:"id"`
	AlarmID    string    `json:"alarm_id"`
	Action     string    `json:"action"` // e.g., "ACKNOWLEDGED", "RESOLVED"
	OperatorID string    `json:"operator_id,omitempty"`
	Notes      string    `json:"notes,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}