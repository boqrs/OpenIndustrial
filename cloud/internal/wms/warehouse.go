package wms

// Warehouse represents a physical building or area for storing goods.
// It is an extension of a Resource of type 'warehouse'.
type Warehouse struct {
	ID         string
	ResourceID string // Foreign key to the resource.Resource
	Name       string
}