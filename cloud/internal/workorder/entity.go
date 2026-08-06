package workorder

import (
	"time"

	"github.com/google/uuid"
)

// WorkOrder represents a command to produce a certain quantity of a specific product.
type WorkOrder struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	ProductID   uuid.UUID `json:"product_id"` // What to produce
	Quantity    int       `json:"quantity"`   // How many to produce
	Status      string    `json:"status"`     // e.g., "Pending", "In_Progress", "Completed", "Cancelled"
	ScheduledAt time.Time `json:"scheduled_at"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
		Status:      "Pending",
		ScheduledAt: scheduledAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}