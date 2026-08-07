package mes

import (
	"time"

	"github.com/google/uuid"
)

// Task represents a single step in the manufacturing process for a product instance.
type Task struct {
	ID                uuid.UUID `json:"id"`
	ProductInstanceID string    `json:"product_instance_id"`
	StationID         string    `json:"station_id"` // The manufacturing station where this task is performed.
	Status            string    `json:"status"`     // e.g., "pending", "in_progress", "completed", "failed"
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}