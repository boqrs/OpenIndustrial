package wms

import "time"

// TransactionType defines the type of inventory movement.
type TransactionType string

const (
	TransactionReceive  TransactionType = "receive"  // Inbound from supplier or production
	TransactionIssue    TransactionType = "issue"    // Outbound to production
	TransactionConsume  TransactionType = "consume"  // Material consumed by a work order
	TransactionTransfer TransactionType = "transfer" // Movement between locations
	TransactionAdjust   TransactionType = "adjust"   // Stock-taking adjustment
	TransactionReturn   TransactionType = "return"   // Return from production
)

// InventoryTransaction is an immutable record of any change in inventory.
// It ensures that all stock movements are event-sourced and traceable.
type InventoryTransaction struct {
	ID              string          `json:"id"`
	Type            TransactionType `json:"type"`
	MaterialBatchID string          `json:"material_batch_id"`
	Quantity        float64         `json:"quantity"` // Can be negative for outbound transactions
	Unit            string          `json:"unit"`
	FromLocationID  string          `json:"from_location_id,omitempty"`
	ToLocationID    string          `json:"to_location_id,omitempty"`
	// ReferenceID links the transaction to a source document, e.g., WorkOrderID, PurchaseOrderID.
	ReferenceID     string          `json:"reference_id,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}