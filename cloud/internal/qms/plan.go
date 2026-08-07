package qms

// InspectionPlan links a quality standard to a specific node in a process definition.
// It specifies that an inspection must occur at this step.
type InspectionPlan struct {
	ID          string `json:"id"`
	StandardID  string `json:"standard_id"`
	// ProcessNodeID is the ID of the process node where the inspection takes place.
	ProcessNodeID string `json:"process_node_id"`
	Name        string `json:"name"`
}

// InspectionItem defines a single metric or attribute to be checked as part of an InspectionPlan.
// It includes the acceptance criteria (e.g., min/max values).
type InspectionItem struct {
	ID       string `json:"id"`
	PlanID   string `json:"plan_id"`
	Name     string `json:"name"`     // e.g., "motor_current", "case_temperature"
	DataType string `json:"data_type"` // e.g., "float", "boolean", "string"
	// Acceptance criteria
	MinValue *float64 `json:"min_value,omitempty"`
	MaxValue *float64 `json:"max_value,omitempty"`
	ExpectedValue *string `json:"expected_value,omitempty"`
}