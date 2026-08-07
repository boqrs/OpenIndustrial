package workorder

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidQuantity is returned when the quantity for a work order is not valid.
	ErrInvalidQuantity = errors.New("invalid quantity")
)

// WorkOrder defines the core entity for a work order.
type WorkOrder struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	ProductID   string    `json:"product_id"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"` // e.g., "Pending", "IN_PROGRESS", "Completed"
	ScheduledAt time.Time `json:"scheduled_at"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StationTask represents a specific task for a product instance at a station,
// derived from a work order.
type StationTask struct {
	ID                string
	WorkOrderID       string
	ProductInstanceID string
	SN                string
	StationID         string
	Status            string // e.g., "pending", "in_progress", "finished"
}

// NewWorkOrder creates a new WorkOrder entity.
func NewWorkOrder(orgID, productID string, quantity int, scheduledAt time.Time) (*WorkOrder, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	now := time.Now().UTC()
	return &WorkOrder{
		ID:          uuid.NewString(),
		OrgID:       orgID,
		ProductID:   productID,
		Quantity:    quantity,
		Status:      "Pending",
		ScheduledAt: scheduledAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}