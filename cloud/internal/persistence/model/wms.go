package model

import (
	"time"

	"github.com/google/uuid"
)

type InventoryStatus string

const (
	InventoryStatusInStock  InventoryStatus = "in_stock"
	InventoryStatusReserved InventoryStatus = "reserved"
	InventoryStatusShipped  InventoryStatus = "shipped"
)

func (s InventoryStatus) String() string {
	return string(s)
}

func (s InventoryStatus) IsValid() bool {
	switch s {
	case InventoryStatusInStock,
		InventoryStatusReserved,
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
//
// Warehouse is a tenant-owned root entity.
// Code is unique inside a tenant, not globally.
type Warehouse struct {
	ID uint `gorm:"primaryKey"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	Code string `gorm:"type:varchar(100);not null"`
	Name string `gorm:"type:varchar(255);not null"`

	Address string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// WarehouseLocation represents a physical location inside a warehouse.
//
// Tenant ownership is inherited through Warehouse.
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
// Tenant ownership is inherited through:
//
//	Device
//	  -> Resource
//	  -> TenantID
//
// A Device has one current inventory record.
//
// Inventory lifecycle:
//
//	in_stock -> reserved -> shipped
//
// A returned shipped device may re-enter:
//
//	shipped -> in_stock
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
// Shipment is a tenant-owned root entity.
//
// SalesOrderID is optional because WMS can also handle
// shipments that do not originate from an internal SalesOrder.
type Shipment struct {
	ID uint `gorm:"primaryKey"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	SalesOrderID *uint `gorm:"index"`

	ExternalOrderID string `gorm:"type:varchar(255);index"`

	Carrier        string `gorm:"type:varchar(100);not null"`
	TrackingNumber string `gorm:"type:varchar(255);index"`

	Status ShipmentStatus `gorm:"type:varchar(50);not null;index"`

	ShippedAt   *time.Time
	DeliveredAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ShipmentItem connects a Device to a Shipment.
//
// SalesOrderItemID is optional for non-SalesOrder shipments.
//
// DeviceID is intentionally NOT globally unique.
//
// This allows:
//
//	shipped
//	  -> returned
//	  -> restocked
//	  -> reshipped
type ShipmentItem struct {
	ID uint `gorm:"primaryKey"`

	ShipmentID uint `gorm:"not null;index"`

	DeviceID uint `gorm:"not null;index"`

	SalesOrderItemID *uint `gorm:"index"`

	CreatedAt time.Time
}

// ShipmentTrackingEvent stores third-party logistics events.
//
// Tenant ownership is inherited through Shipment.
type ShipmentTrackingEvent struct {
	ID uint `gorm:"primaryKey"`

	ShipmentID uint `gorm:"not null;index"`

	ExternalEventID string `gorm:"type:varchar(255);not null;index"`

	Status ShipmentStatus `gorm:"type:varchar(50);not null;index"`

	OccurredAt time.Time `gorm:"not null;index"`

	Location    string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:text"`

	CreatedAt time.Time
}
