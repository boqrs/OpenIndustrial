package mes

import "time"

// ProductionLine represents a physical or logical assembly line within a factory.
// It is an extension of a Resource of type 'line'.
type ProductionLine struct {
	ID         string
	ResourceID string // Foreign key to the resource.Resource
	Name       string
	Code       string
	Status     LineStatus
	CreatedAt  time.Time
}

// LineStatus defines the operational status of a production line.
type LineStatus string

const (
	LineRunning     LineStatus = "running"
	LineStopped     LineStatus = "stopped"
	LineMaintenance LineStatus = "maintenance"
)