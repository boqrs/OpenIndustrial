package workorder

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Constants for WorkOrder and Task statuses
const (
	StatusPending     = "Pending"
	StatusInProgress  = "InProgress"
	StatusCompleted   = "Completed"
	StatusTaskPending = "Pending"
)

var (
	// ErrInvalidQuantity is returned when the quantity for a work order is not valid.
	ErrInvalidQuantity = errors.New("invalid quantity")
)

// WorkOrder defines the core entity for a work order.
type WorkOrder struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"` // e.g., "Pending", "InProgress", "Completed"
	ScheduledAt time.Time `json:"scheduled_at"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StationTask represents a specific task for a product instance at a station,
// derived from a work order.
type StationTask struct {
	ID                uuid.UUID `json:"id"`
	WorkOrderID       uuid.UUID `json:"work_order_id"`
	ProductInstanceID uuid.UUID `json:"product_instance_id"`
	SN                string    `json:"sn"`
	StationID         uuid.UUID `json:"station_id"`
	Status            string    `json:"status"` // e.g., "pending", "in_progress", "finished"
	Result            string    `json:"result,omitempty"`
	CompletedAt       time.Time `json:"completed_at,omitempty"`
}

// NewWorkOrder creates a new WorkOrder entity.
func NewWorkOrder(orgID, productID uuid.UUID, quantity int, scheduledAt time.Time) (*WorkOrder, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	now := time.Now().UTC()
	return &WorkOrder{
		ID:          uuid.New(),
		OrgID:       orgID,
		ProductID:   productID,
		Quantity:    quantity,
		Status:      StatusPending,
		ScheduledAt: scheduledAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}