package telemetry

// Metric represents a definition of an industrial indicator, like 'motor.current'.
type Metric struct {
	ID       string
	Code     string // e.g., "motor.current", "temperature"
	Name     string
	Unit     string // e.g., "A", "℃"
	DataType string // e.g., "float", "integer", "boolean"
}

// DeviceMetric binds a specific metric to a device and its address in the physical world.
type DeviceMetric struct {
	ID       string
	DeviceID string // Foreign key to the device
	MetricID string // Foreign key to Metric
	Address  string // The address on the device (e.g., "DB10.20", "40001")
}