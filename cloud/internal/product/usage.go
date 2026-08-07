package product

import "time"

// MaterialUsageRecord links a finished product instance to the specific material batches consumed during its production.
// This is the cornerstone of end-to-end traceability.
type MaterialUsageRecord struct {
	ID                string
	ProductInstanceID string    // The final product SN
	MaterialBatchID   string    // The batch of material used
	Quantity          float64   // The quantity of material from that batch
	Unit              string
	StationID         string    // Where the material was consumed
	Timestamp         time.Time // When it was consumed
}