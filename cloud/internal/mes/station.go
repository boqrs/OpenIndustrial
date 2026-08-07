package mes

// Station represents a physical location or a piece of equipment on a production line.
// It is an extension of a Resource of type 'station'.
type Station struct {
	ID         string
	ResourceID string // Foreign key to the resource.Resource
	LineID     string // Foreign key to ProductionLine
	Name       string
	Type       string // e.g., "manual_assembly", "automated_test", "packing"
	Status     string // e.g., "idle", "running", "error"
}

// Capability defines a specific function that a station can perform.
// This decouples the process from the physical station.
type Capability struct {
	ID        string
	StationID string // Foreign key to Station
	Name      string // e.g., "firmware_upgrade", "pressure_test"
}