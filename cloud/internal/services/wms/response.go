package wms

import "time"

type WarehouseResponse struct {
	ID uint `json:"id"`

	Code    string `json:"code"`
	Name    string `json:"name"`
	Address string `json:"address"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LocationResponse struct {
	ID uint `json:"id"`

	WarehouseID uint `json:"warehouse_id"`

	Code string `json:"code"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InventoryResponse struct {
	ID uint `json:"id"`

	DeviceID uint `json:"device_id"`

	WarehouseID uint `json:"warehouse_id"`
	LocationID  uint `json:"location_id"`

	Status string `json:"status"`

	InboundAt  time.Time  `json:"inbound_at"`
	OutboundAt *time.Time `json:"outbound_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShipmentResponse struct {
	ID uint `json:"id"`

	ExternalOrderID string `json:"external_order_id"`

	Carrier        string `json:"carrier"`
	TrackingNumber string `json:"tracking_number"`

	Status string `json:"status"`

	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`

	Items []*ShipmentItemResponse `json:"items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShipmentItemResponse struct {
	ID uint `json:"id"`

	ShipmentID uint `json:"shipment_id"`
	DeviceID   uint `json:"device_id"`

	CreatedAt time.Time `json:"created_at"`
}

type TrackingEventResponse struct {
	ID uint `json:"id"`

	ShipmentID uint `json:"shipment_id"`

	ExternalEventID string `json:"external_event_id"`

	Status string `json:"status"`

	OccurredAt time.Time `json:"occurred_at"`

	Location    string `json:"location"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
}