package product

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

// EventHandler handles incoming events related to products.
type EventHandler struct {
	service *Service
}

// NewEventHandler creates a new event handler.
func NewEventHandler(service *Service) *EventHandler {
	return &EventHandler{service: service}
}

// ProductEventMessage defines the structure of an incoming product event.
type ProductEventMessage struct {
	ProductInstanceID string                 `json:"product_instance_id"`
	EventType         string                 `json:"event_type"`
	Timestamp         int64                  `json:"timestamp"` // Unix nano
	Location          string                 `json:"location"`
	Properties        map[string]interface{} `json:"properties"`
}

// HandleProductEvent processes a raw event message.
func (h *EventHandler) HandleProductEvent(ctx context.Context, rawMessage []byte) error {
	var msg ProductEventMessage
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		log.Printf("Failed to unmarshal product event: %v", err)
		return err
	}

	instanceID, err := uuid.Parse(msg.ProductInstanceID)
	if err != nil {
		log.Printf("Invalid ProductInstanceID in event message: %v", err)
		return err
	}

	event := &LifecycleEvent{
		ProductInstanceID: instanceID,
		Type:              msg.EventType,
		Timestamp:         time.Unix(0, msg.Timestamp),
		Location:          msg.Location,
		Properties:        msg.Properties,
	}

	return h.service.AddLifecycleEvent(ctx, event)
}