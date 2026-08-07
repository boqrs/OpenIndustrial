package telemetry

import (
	"time"
)

// DataPoint represents a single metric value for a given resource at a specific time.
// This is the fundamental unit of time-series data in the new resource model.
type DataPoint struct {
	ResourceID string    `json:"resource_id"`
	Metric     string    `json:"metric"` // e.g., "temperature", "energy_consumption"
	Value      float64   `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
}