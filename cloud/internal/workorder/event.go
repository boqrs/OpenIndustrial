package workorder

import (
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/event"
)

const (
	EventProductionFinished = "production.finished"
)

// ProductionFinishedEvent is emitted when one product instance
// has finished its production process within a work order.
type ProductionFinishedEvent struct {
	WorkOrderID       string
	ProductInstanceID string
	SN                string
	StationID         string
	Result            string // e.g., "pass", "fail"
	FinishedAt        time.Time
}

// ToDomainEvent converts the business-specific event into a generic domain event.
func (e ProductionFinishedEvent) ToDomainEvent() event.Event {
	return event.Event{
		// ID should be set by the event bus or publisher if needed
		Name:        EventProductionFinished,
		Aggregate:   "WorkOrder",
		AggregateID: e.WorkOrderID,
		Source:      "mes",
		Data:        e, // The entire struct is the payload
		CreatedAt:   time.Now().UTC(),
	}
}