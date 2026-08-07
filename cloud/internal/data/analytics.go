package data

import "time"

// Feature represents a derived, analytical data point that can be used for
// higher-level analysis, such as AI/ML models. It is often the output of
// an aggregation or transformation function.
type Feature struct {
	ResourceID string    `json:"resource_id"`
	// Name is the identifier of the feature, e.g., "vibration_fft_peak", "current_7day_avg".
	Name       string    `json:"name"`
	Value      float64   `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
}