package wms

import "time"

// Inventory represents the quantity of a specific material batch at a specific location.
// A location is a Resource of type 'location'.
type Inventory struct {
	ID         string    `json:"id"`
	// LocationID is the ID of the Resource representing the warehouse location.
	LocationID string    `json:"location_id"`
	BatchID    string    `json:"batch_id"`
	Quantity   float64   `json:"quantity"`
	Unit       string    `json:"unit"`
	UpdatedAt  time.Time `json:"updated_at"`
}