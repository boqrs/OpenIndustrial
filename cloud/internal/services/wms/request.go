package wms

import "time"

type CreateWarehouseRequest struct {
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

type CreateLocationRequest struct {
	WarehouseID uint `json:"warehouse_id" binding:"required"`

	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type StockInRequest struct {
	DeviceID uint `json:"device_id" binding:"required"`

	WarehouseID uint `json:"warehouse_id" binding:"required"`
	LocationID  uint `json:"location_id" binding:"required"`

	InboundAt *time.Time `json:"inbound_at"`
}

type CreateShipmentRequest struct {
	// SalesOrderID links this shipment to an internal SalesOrder.
	//
	// It is optional because WMS may also process shipments
	// that are not generated from an internal SalesOrder.
	SalesOrderID *uint `json:"sales_order_id"`

	// ExternalOrderID is the order identifier from an external
	// ERP / e-commerce / logistics system.
	ExternalOrderID string `json:"external_order_id"`

	Carrier        string `json:"carrier" binding:"required"`
	TrackingNumber string `json:"tracking_number" binding:"required"`

	// Items is required for SalesOrder-linked shipments because
	// every device must be mapped to a SalesOrderItem.
	Items []CreateShipmentItemRequest `json:"items"`

	// DeviceIDs is retained for backward compatibility for
	// non-SalesOrder shipments.
	//
	// When SalesOrderID is provided, Items must be used.
	DeviceIDs []uint `json:"device_ids"`
}

type CreateShipmentItemRequest struct {
	DeviceID uint `json:"device_id" binding:"required"`

	// SalesOrderItemID is required when SalesOrderID is provided.
	SalesOrderItemID *uint `json:"sales_order_item_id"`
}

type TrackingEventRequest struct {
	ExternalEventID string `json:"external_event_id" binding:"required"`

	Status ShipmentStatusRequest `json:"status" binding:"required"`

	OccurredAt *time.Time `json:"occurred_at"`

	Location    string `json:"location"`
	Description string `json:"description"`
}

type ShipmentStatusRequest string

const (
	ShipmentStatusRequestInTransit      ShipmentStatusRequest = "in_transit"
	ShipmentStatusRequestOutForDelivery ShipmentStatusRequest = "out_for_delivery"
	ShipmentStatusRequestDelivered      ShipmentStatusRequest = "delivered"
	ShipmentStatusRequestException      ShipmentStatusRequest = "exception"
)

func (s ShipmentStatusRequest) ToModel() (string, bool) {
	switch s {
	case ShipmentStatusRequestInTransit,
		ShipmentStatusRequestOutForDelivery,
		ShipmentStatusRequestDelivered,
		ShipmentStatusRequestException:
		return string(s), true
	default:
		return "", false
	}
}