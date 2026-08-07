package logistics

import (
	"encoding/json"
	"time"
)

// DeliveryEventType defines the type of a logistics event.
type DeliveryEventType string

const (
	EventPickedUp      DeliveryEventType = "picked_up"
	EventInTransit     DeliveryEventType = "in_transit"
	EventArrivedAtHub  DeliveryEventType = "arrived_at_hub"
	EventOutForDelivery DeliveryEventType = "out_for_delivery"
	EventDelivered     DeliveryEventType = "delivered"
	EventFailedAttempt DeliveryEventType = "failed_attempt"
)

// DeliveryEvent is an immutable record of a logistics activity.
type DeliveryEvent struct {
	ID         string
	ShipmentID string
	Type       DeliveryEventType
	Location   string
	Notes      string
	Data       json.RawMessage
	Timestamp  time.Time
}