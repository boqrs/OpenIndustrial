package qms

import "time"

// DefectSeverity defines the criticality of a defect.
type DefectSeverity string

const (
	SeverityCritical DefectSeverity = "Critical"
	SeverityMajor    DefectSeverity = "Major"
	SeverityMinor    DefectSeverity = "Minor"
)

// Defect represents a specific type of non-conformance found during inspection.
type Defect struct {
	ID          string         `json:"id"`
	// ProductID links the defect to the specific product instance.
	ProductID   string         `json:"product_id"`
	// InspectionRecordID links the defect to the inspection where it was found.
	InspectionRecordID string `json:"inspection_record_id"`
	// Code is a standardized identifier for the defect type, e.g., "MOTOR_NOISE_HIGH".
	Code        string         `json:"code"`
	Description string         `json:"description,omitempty"`
	Severity    DefectSeverity `json:"severity"`
	ReportedAt  time.Time      `json:"reported_at"`
}