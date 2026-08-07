package telemetry

import "time"

// TelemetryContext provides the business context for a stream of telemetry data.
// It links data points to a specific product, station, and task during production.
type TelemetryContext struct {
	ID                string
	DeviceID          string
	ProductInstanceID string
	StationID         string
	TaskID            string
	StartAt           time.Time
	EndAt             *time.Time
}