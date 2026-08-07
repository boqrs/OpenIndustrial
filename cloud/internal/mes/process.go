package mes

// ProcessDefinition is the blueprint for manufacturing a specific product model.
// It contains a versioned set of steps.
type ProcessDefinition struct {
	ID             string
	ProductModelID string // Foreign key to product.ProductModel
	Name           string
	Version        int
	Status         string // e.g., "draft", "active", "archived"
}

// ProcessStep is a single, ordered step within a ProcessDefinition.
type ProcessStep struct {
	ID                string
	ProcessID         string // Foreign key to ProcessDefinition
	Name              string
	Sequence          int
	StationCapability string // The required capability for a station to perform this step.
	Required          bool
}