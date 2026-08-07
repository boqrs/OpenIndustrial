package product

import "time"

// WorkOrder represents a production order to manufacture a certain quantity of a product model.
type WorkOrder struct {
	ID             string
	OrgID          string
	OrderNo        string
	ProductModelID string
	Quantity       int
	Status         string // e.g., "pending", "in_progress", "completed"
	CreatedAt      time.Time
}

// WorkOrderItem links a specific ProductInstance (SN) to a WorkOrder.
type WorkOrderItem struct {
	ID                string
	WorkOrderID       string
	ProductInstanceID string
}

// ProcessDefinition defines a step in the manufacturing process for a product model.
type ProcessDefinition struct {
	ID             string
	ProductModelID string
	Name           string // e.g., "Assembly", "Firmware Flashing", "Final Test"
	Sequence       int    // The order of this process step
}