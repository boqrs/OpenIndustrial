package qms

import (
	"encoding/json"
	"time"
)

// InspectionResultStatus defines the outcome of an inspection.
type InspectionResultStatus string

const (
	ResultPass   InspectionResultStatus = "PASS"
	ResultFail   InspectionResultStatus = "FAIL"
	ResultWarn   InspectionResultStatus = "WARN"
	ResultSkip   InspectionResultStatus = "SKIPPED"
)

// InspectionRecord is an immutable log of an inspection performed on a product.
// It serves as a verifiable quality certificate for a specific production step.
type InspectionRecord struct {
	ID          string                 `json:"id"`
	// ProductID links the record to the specific product instance being inspected.
	ProductID   string                 `json:"product_id"`
	PlanID      string                 `json:"plan_id"`
	// ResourceID is the station or equipment where the inspection was performed.
	ResourceID  string                 `json:"resource_id"`
	OperatorID  string                 `json:"operator_id,omitempty"`
	Result      InspectionResultStatus `json:"result"`
	// Data contains the raw measurement values collected during the inspection.
	// e.g., {"motor_current": 1.5, "temperature": 45.2}
	Data        json.RawMessage        `json:"data"`
	CreatedAt   time.Time              `json:"created_at"`
}