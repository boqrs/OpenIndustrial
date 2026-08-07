package logistics

import "time"

// ShipmentStatus defines the lifecycle of a shipment.
type ShipmentStatus string

const (
	ShipmentPending   ShipmentStatus = "pending"
	ShipmentReady     ShipmentStatus = "ready_for_pickup"
	ShipmentInTransit ShipmentStatus = "in_transit"
	ShipmentDelivered ShipmentStatus = "delivered"
	ShipmentCancelled ShipmentStatus = "cancelled"
)

// Shipment represents a collection of products being transported.
type Shipment struct {
	ID              string
	OrderID         string // Reference to a sales order
	Status          ShipmentStatus
	Carrier         string // Name of the logistics company
	TrackingNumber  string
	EstimatedPickup *time.Time
	EstimatedDelivery *time.Time
	CreatedAt       time.Time
}

// ShipmentItem links a specific product instance to a shipment.
type ShipmentItem struct {
	ID                string
	ShipmentID        string // Foreign key to Shipment
	ProductInstanceID string // Foreign key to product.ProductInstance
}