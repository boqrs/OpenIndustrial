package telemetry

import (
	"context"
	"time"
)

// Repository defines the interface for storing and retrieving telemetry data.
// This interface abstracts away the underlying storage (e.g., TimescaleDB, Redis).
type Repository interface {
	// Time-series data methods
	SaveDataPoints(ctx context.Context, points []*DataPoint) error
	GetDataPoints(ctx context.Context, deviceID, metricID string, start, end time.Time) ([]*DataPoint, error)

	// Real-time state methods
	UpdateDeviceState(ctx context.Context, state *DeviceState) error
	GetDeviceState(ctx context.Context, deviceID string) (*DeviceState, error)

	// Context and Binding methods
	CreateContext(ctx context.Context, context *TelemetryContext) error
	GetActiveContext(ctx context.Context, deviceID string) (*TelemetryContext, error)

	// Alarm and Rule methods
	CreateAlarm(ctx context.Context, alarm *Alarm) error
	GetRulesForMetric(ctx context.Context, metricID string) ([]*Rule, error)
}