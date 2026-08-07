package telemetry

// Rule defines a condition that, when met, triggers an alarm.
type Rule struct {
	ID         string
	MetricID   string
	Expression string // e.g., "value > 80", "value < 10 && value > 5"
	Duration   int    // How long the condition must be true (in seconds) before alarming.
	Level      AlarmLevel
	Message    string
}