package wms

// Material defines a specific type of raw material, component, or substance
// used in the manufacturing process. It is the master data for inventory.
type Material struct {
	ID    string `json:"id"`
	OrgID string `json:"org_id"`
	// Code is the unique identifier for the material (e.g., SKU).
	Code  string `json:"code"`
	Name  string `json:"name"`
	// Unit of measure for this material (e.g., "pcs", "kg", "m").
	Unit  string `json:"unit"`
}