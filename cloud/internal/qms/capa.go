package qms

import "time"

// CAPAStatus defines the lifecycle of a corrective/preventive action.
type CAPAStatus string

const (
	CAPAStatusOpen        CAPAStatus = "Open"
	CAPAStatusInProgress  CAPAStatus = "In Progress"
	CAPAStatusPendingReview CAPAStatus = "Pending Review"
	CAPAStatusClosed      CAPAStatus = "Closed"
	CAPAStatusCancelled   CAPAStatus = "Cancelled"
)

// CAPA (Corrective and Preventive Action) is a formal process to investigate,
// address, and prevent the recurrence of defects or quality issues.
type CAPA struct {
	ID       string   `json:"id"`
	// DefectIDs links the CAPA to one or more defects it addresses.
	DefectIDs []string `json:"defect_ids"`
	Title    string   `json:"title"`
	// RootCauseAnalysis is a detailed description of the investigation findings.
	RootCauseAnalysis string `json:"root_cause_analysis"`
	// CorrectiveAction describes the immediate steps taken to fix the issue.
	CorrectiveAction string `json:"corrective_action"`
	// PreventiveAction describes the long-term steps to prevent recurrence.
	PreventiveAction string `json:"preventive_action"`
	Status   CAPAStatus `json:"status"`
	OwnerID  string   `json:"owner_id"` // User responsible for the CAPA
	DueDate  *time.Time `json:"due_date,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ClosedAt *time.Time `json:"closed_at,omitempty"`
}