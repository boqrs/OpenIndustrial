package wms

import (
	"context"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}

type Repository interface {
	// Warehouse
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	GetWarehouseByID(ctx context.Context, id uint) (*model.Warehouse, error)

	// Location
	CreateLocation(ctx context.Context, location *model.WarehouseLocation) error
	GetLocationByID(ctx context.Context, id uint) (*model.WarehouseLocation, error)

	// Inventory
	GetInventoryByDeviceID(
		ctx context.Context,
		deviceID uint,
	) (*model.DeviceInventory, error)

	GetInventoryByDeviceIDForUpdateTx(
		ctx context.Context,
		deviceID uint,
	) (*model.DeviceInventory, error)

	CreateInventoryTx(
		ctx context.Context,
		inventory *model.DeviceInventory,
	) error

	UpdateInventoryTx(
		ctx context.Context,
		inventory *model.DeviceInventory,
	) error

	// Shipment
	CreateShipmentTx(
		ctx context.Context,
		shipment *model.Shipment,
	) error

	CreateShipmentItemsTx(
		ctx context.Context,
		items []*model.ShipmentItem,
	) error

	GetShipmentByID(
		ctx context.Context,
		id uint,
	) (*model.Shipment, error)

	GetShipmentByIDForUpdateTx(
		ctx context.Context,
		id uint,
	) (*model.Shipment, error)

	ListShipmentItems(
		ctx context.Context,
		shipmentID uint,
	) ([]*model.ShipmentItem, error)

	UpdateShipmentTx(
		ctx context.Context,
		shipment *model.Shipment,
	) error

	// Tracking
	GetTrackingEventByExternalID(
		ctx context.Context,
		shipmentID uint,
		externalEventID string,
	) (*model.ShipmentTrackingEvent, error)

	CreateTrackingEventTx(
		ctx context.Context,
		event *model.ShipmentTrackingEvent,
	) error

	ListTrackingEvents(
		ctx context.Context,
		shipmentID uint,
	) ([]*model.ShipmentTrackingEvent, error)
}

type Service interface {
	CreateWarehouse(
		ctx context.Context,
		req *CreateWarehouseRequest,
	) (*WarehouseResponse, error)

	GetWarehouse(
		ctx context.Context,
		id uint,
	) (*WarehouseResponse, error)

	CreateLocation(
		ctx context.Context,
		req *CreateLocationRequest,
	) (*LocationResponse, error)

	GetDeviceInventory(
		ctx context.Context,
		deviceID uint,
	) (*InventoryResponse, error)

	StockIn(
		ctx context.Context,
		req *StockInRequest,
	) (*InventoryResponse, error)

	CreateShipment(
		ctx context.Context,
		req *CreateShipmentRequest,
	) (*ShipmentResponse, error)

	StockOut(
		ctx context.Context,
		shipmentID uint,
	) error

	GetShipment(
		ctx context.Context,
		shipmentID uint,
	) (*ShipmentResponse, error)

	AddTrackingEvent(
		ctx context.Context,
		shipmentID uint,
		req *TrackingEventRequest,
	) error

	ListTrackingEvents(
		ctx context.Context,
		shipmentID uint,
	) ([]*TrackingEventResponse, error)
}
