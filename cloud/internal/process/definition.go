package process

// ProcessDefinition is a template for a production or business process.
// It defines the name, version, and the product it applies to.
type ProcessDefinition struct {
	ID          string `json:"id"`
	OrgID       string `json:"org_id"`
	// ProductType links the definition to a specific type of product resource.
	ProductType string `json:"product_type"`
	Name        string `json:"name"`
	Version     int    `json:"version"`
	// Status indicates if the definition is active, archived, or in draft.
	Status      string `json:"status"` // e.g., "DRAFT", "ACTIVE", "ARCHIVED"
}