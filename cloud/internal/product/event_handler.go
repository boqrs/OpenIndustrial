package product

import (
	"context"
	"fmt"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/event"
)

const (
	ProductionFinishedEventName = "production.finished"
)

// ProductionFinishedEvent is the event payload for when a product instance is finished.
type ProductionFinishedEvent struct {
	WorkOrderID       string
	ProductInstanceID string
	SN                string
	StationID         string
	Result            string
	FinishedAt        time.Time
}

// ToDomainEvent converts the specific event to a generic event.Event for the bus.
func (pfe ProductionFinishedEvent) ToDomainEvent() event.Event {
	return event.Event{
		Name: ProductionFinishedEventName,
		Data: pfe,
	}
}

// EventHandler handles events related to the product domain.
type EventHandler struct {
	service Service
}

// NewEventHandler creates a new event handler for the product domain.
func NewEventHandler(service Service) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

// Register handles the subscription of relevant events.
func (h *EventHandler) Register(bus event.Bus) {
	bus.Subscribe(ProductionFinishedEventName, h)
}

// Handle is the method that satisfies the event.Handler interface.
func (h *EventHandler) Handle(e event.Event) error {
	switch e.Name {
	case ProductionFinishedEventName:
		return h.handleProductionFinished(e)
	}
	return nil
}

// handleProductionFinished processes the production finished event.
func (h *EventHandler) handleProductionFinished(e event.Event) error {
	payload, ok := e.Data.(ProductionFinishedEvent)
	if !ok {
		err := fmt.Errorf("invalid payload type for event %s", ProductionFinishedEventName)
		fmt.Println(err)
		return err
	}

	lifecycle := LifecycleEvent{
		ProductInstanceID: payload.ProductInstanceID,
		Type:              "MANUFACTURED",
		Timestamp:         payload.FinishedAt,
		Metadata: map[string]interface{}{
			"work_order_id": payload.WorkOrderID,
			"station_id":    payload.StationID,
			"sn":            payload.SN,
			"result":        payload.Result,
		},
	}

	if err := h.service.AppendLifecycle(context.Background(), lifecycle); err != nil {
		fmt.Printf("Error appending lifecycle event: %v\n", err)
		return err
	}
	return nil
}