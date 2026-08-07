package data

// MetricPoint is the standardized, unified model for any single piece of time-series data
// ingested into the platform. It enriches a simple value with critical context.
type MetricPoint struct {
	// ResourceID links the metric to a specific entity in the Digital Twin Graph.
	ResourceID string `json:"resource_id"`
	// ProductID optionally links the metric to a product instance, crucial for lifecycle analysis.
	ProductID string `json:"product_id,omitempty"`
	// Metric is the name of the measurement, e.g., "motor_current".
	Metric    string  `json:"metric"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit,omitempty"`
	// Timestamp is a Unix timestamp in milliseconds.
	Timestamp int64   `json:"timestamp"`
}

// MetricDefinition defines a valid metric that can be ingested for a specific resource type.
// This acts as a schema and validation rule for incoming telemetry.
type MetricDefinition struct {
	ID           string `json:"id"`
	// ResourceType specifies which type of resource this metric applies to (e.g., "robot", "plc").
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
	DataType     string `json:"data_type"` // e.g., "float", "integer", "boolean"
	Unit         string `json:"unit"`
	Description  string `json:"description,omitempty"`
}