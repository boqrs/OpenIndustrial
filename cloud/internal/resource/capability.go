package resource

// Capability represents an action or function that a resource can perform.
// This defines the "methods" of a digital twin.
type Capability struct {
	// ResourceID links the capability to its parent resource.
	ResourceID string `json:"resource_id"`
	// Code is a unique identifier for the capability, e.g., "motion.control", "modbus.read".
	Code string `json:"code"`
	// Version of the capability's implementation.
	Version string `json:"version,omitempty"`
}