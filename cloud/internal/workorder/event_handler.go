package workorder

import (
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/event"
)

// EventHandler for workorder domain.
// Currently, it does not subscribe to any events, but it's good practice
// to have the structure ready.
type EventHandler struct {
	service Service
}

// NewEventHandler creates a new workorder event handler.
func NewEventHandler(service Service) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

// Register subscribes to events.
func (h *EventHandler) Register(bus event.Bus) {
	// No subscriptions needed for now.
}

// Handle processes events. This method makes the handler satisfy the event.Handler interface.
func (h *EventHandler) Handle(e event.Event) error {
	// No events to handle yet.
	return nil
}