package model

import "time"

type InventoryStatus string

const (
	InventoryStatusInStock InventoryStatus = "in_stock"
	InventoryStatusShipped InventoryStatus = "shipped"
)

func (s InventoryStatus) String() string {
	return string(s)
}

func (s InventoryStatus) IsValid() bool {
	switch s {
	case InventoryStatusInStock,
		InventoryStatusShipped:
		return true
	default:
		return false
	}
}

type ShipmentStatus string

const (
	ShipmentStatusCreated        ShipmentStatus = "created"
	ShipmentStatusInTransit      ShipmentStatus = "in_transit"
	ShipmentStatusOutForDelivery ShipmentStatus = "out_for_delivery"
	ShipmentStatusDelivered      ShipmentStatus = "delivered"
	ShipmentStatusException      ShipmentStatus = "exception"
	ShipmentStatusCancelled      ShipmentStatus = "cancelled"
)

func (s ShipmentStatus) String() string {
	return string(s)
}

func (s ShipmentStatus) IsValid() bool {
	switch s {
	case ShipmentStatusCreated,
		ShipmentStatusInTransit,
		ShipmentStatusOutForDelivery,
		ShipmentStatusDelivered,
		ShipmentStatusException,
		ShipmentStatusCancelled:
		return true
	default:
		return false
	}
}

// Warehouse represents a physical warehouse.
type Warehouse struct {
	ID uint `gorm:"primaryKey"`

	Code string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Name string `gorm:"type:varchar(255);not null"`

	Address string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// WarehouseLocation represents a physical location inside a warehouse.
type WarehouseLocation struct {
	ID uint `gorm:"primaryKey"`

	WarehouseID uint `gorm:"not null;index"`

	Code string `gorm:"type:varchar(100);not null"`
	Name string `gorm:"type:varchar(255);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// DeviceInventory represents the current WMS inventory state of a Device.
//
// A Device has one current inventory record.
// Re-stock after a return updates the same record.
type DeviceInventory struct {
	ID uint `gorm:"primaryKey"`

	DeviceID uint `gorm:"not null;uniqueIndex"`

	WarehouseID uint `gorm:"not null;index"`
	LocationID  uint `gorm:"not null;index"`

	Status InventoryStatus `gorm:"type:varchar(50);not null;index"`

	InboundAt  time.Time
	OutboundAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Shipment represents an outbound shipment.
//
// WMS deliberately does not reference CustomerID or UserID.
// At this stage the external order/channel identity is enough.
type Shipment struct {
	ID uint `gorm:"primaryKey"`

	ExternalOrderID string `gorm:"type:varchar(255);index"`

	Carrier string `gorm:"type:varchar(100);not null"`

	TrackingNumber string `gorm:"type:varchar(255);index"`

	Status ShipmentStatus `gorm:"type:varchar(50);not null;index"`

	ShippedAt   *time.Time
	DeliveredAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ShipmentItem connects a Device to a Shipment.
//
// DeviceID is intentionally NOT globally unique.
// This allows a future return -> restock -> reship lifecycle.
type ShipmentItem struct {
	ID uint `gorm:"primaryKey"`

	ShipmentID uint `gorm:"not null;index"`
	DeviceID   uint `gorm:"not null;index"`

	CreatedAt time.Time
}

// ShipmentTrackingEvent stores third-party logistics events.
type ShipmentTrackingEvent struct {
	ID uint `gorm:"primaryKey"`

	ShipmentID uint `gorm:"not null;index"`

	ExternalEventID string `gorm:"type:varchar(255);index"`

	Status ShipmentStatus `gorm:"type:varchar(50);not null;index"`

	OccurredAt time.Time `gorm:"not null;index"`

	Location    string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:text"`

	CreatedAt time.Time
}
