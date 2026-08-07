package wms

// Location represents a specific place within a warehouse where materials are stored.
// e.g., Aisle A, Rack 01, Shelf 01.
type Location struct {
	ID          string
	WarehouseID string // Foreign key to Warehouse
	Code        string // A human-readable code, e.g., "A01-01-01"
	Name        string
}