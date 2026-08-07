package wms

import "time"

// MaterialBatch represents a specific lot or batch of a material,
// typically from a single supplier and production run. This is crucial for traceability.
type MaterialBatch struct {
	ID         string    `json:"id"`
	MaterialID string    `json:"material_id"`
	BatchNo    string    `json:"batch_no"`
	SupplierID string    `json:"supplier_id,omitempty"`
	Quantity   float64   `json:"quantity"`
	Unit       string    `json:"unit"`
	ProducedAt time.Time `json:"produced_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}