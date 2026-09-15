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
	ExternalOrderID string `json:"external_order_id"`

	Carrier        string `json:"carrier" binding:"required"`
	TrackingNumber string `json:"tracking_number" binding:"required"`

	DeviceIDs []uint `json:"device_ids" binding:"required,min=1"`
}

type TrackingEventRequest struct {
	ExternalEventID string `json:"external_event_id"`

	Status ShipmentStatusRequest `json:"status" binding:"required"`

	OccurredAt *time.Time `json:"occurred_at"`

	Location    string `json:"location"`
	Description string `json:"description"`
}

type ShipmentStatusRequest string

const (
	ShipmentStatusRequestCreated        ShipmentStatusRequest = "created"
	ShipmentStatusRequestInTransit      ShipmentStatusRequest = "in_transit"
	ShipmentStatusRequestOutForDelivery ShipmentStatusRequest = "out_for_delivery"
	ShipmentStatusRequestDelivered      ShipmentStatusRequest = "delivered"
	ShipmentStatusRequestException      ShipmentStatusRequest = "exception"
	ShipmentStatusRequestCancelled      ShipmentStatusRequest = "cancelled"
)

func (s ShipmentStatusRequest) ToModel() (string, bool) {
	switch s {
	case ShipmentStatusRequestCreated,
		ShipmentStatusRequestInTransit,
		ShipmentStatusRequestOutForDelivery,
		ShipmentStatusRequestDelivered,
		ShipmentStatusRequestException,
		ShipmentStatusRequestCancelled:
		return string(s), true
	default:
		return "", false
	}
}
