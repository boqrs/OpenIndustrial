package wms

import "time"

// NodeMaterial defines the material requirements for a specific process node.
// This forms part of the Manufacturing BOM (MBOM).
type NodeMaterial struct {
	// NodeID is the ID of the process node that requires the material.
	NodeID     string  `json:"node_id"`
	// MaterialID is the ID of the required material master data.
	MaterialID string  `json:"material_id"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit"`
}

// MaterialConsumption is an immutable record that links a final product
// to the specific material batch consumed during its production. This is the
// core of the material genealogy.
type MaterialConsumption struct {
	// ProductID is the ID of the final product instance (e.g., by serial number).
	ProductID       string    `json:"product_id"`
	MaterialBatchID string    `json:"material_batch_id"`
	// ProcessNodeID identifies at which step the material was consumed.
	ProcessNodeID   string    `json:"process_node_id"`
	Quantity        float64   `json:"quantity"`
	Unit            string    `json:"unit"`
	ConsumedAt      time.Time `json:"consumed_at"`
}