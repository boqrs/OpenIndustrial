package workorder

import (
	"time"

	"github.com/google/uuid"
)

// WorkOrderCreatedEvent is published when a new work order is created.
type WorkOrderCreatedEvent struct {
	EventID     uuid.UUID `json:"event_id"`
	WorkOrderID uuid.UUID `json:"work_order_id"`
	OrgID       uuid.UUID `json:"org_id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int       `json:"quantity"`
	Timestamp   time.Time `json:"timestamp"`
}

// WorkOrderStatusChangedEvent is published when a work order's status changes.
type WorkOrderStatusChangedEvent struct {
	EventID     uuid.UUID `json:"event_id"`
	WorkOrderID uuid.UUID `json:"work_order_id"`
	OrgID       uuid.UUID `json:"org_id"`
	OldStatus   string    `json:"old_status"`
	NewStatus   string    `json:"new_status"`
	Timestamp   time.Time `json:"timestamp"`
}