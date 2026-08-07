package telemetry

import "time"

// AlarmLevel defines the severity of an alarm.
type AlarmLevel string

const (
	AlarmInfo     AlarmLevel = "info"
	AlarmWarning  AlarmLevel = "warning"
	AlarmCritical AlarmLevel = "critical"
)

// Alarm represents an event triggered when a metric violates a predefined rule.
type Alarm struct {
	ID        string
	DeviceID  string
	MetricID  string
	Level     AlarmLevel
	Message   string
	Timestamp time.Time
	Acked     bool // Has the alarm been acknowledged?
	AckedBy   string
	AckedAt   *time.Time
}