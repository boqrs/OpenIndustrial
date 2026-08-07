package qms

// QualityStandard defines a set of quality criteria and specifications for a
// particular product model or process. It is the master document for quality control.
type QualityStandard struct {
	ID             string `json:"id"`
	OrgID          string `json:"org_id"`
	// ProductModelID links the standard to a specific product model resource.
	ProductModelID string `json:"product_model_id"`
	Name           string `json:"name"`
	Version        int    `json:"version"`
	// Status indicates if the standard is a draft, active, or archived.
	Status         string `json:"status"` // e.g., "DRAFT", "ACTIVE", "ARCHIVED"
}