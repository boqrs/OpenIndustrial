package postgres

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
	"gorm.io/gorm/clause"
)

type WMSRepository struct {
	db *database.DBProvider
}

func NewWMSRepository(db *database.DBProvider) *WMSRepository {
	return &WMSRepository{
		db: db,
	}
}

// -----------------------------------------------------------------------------
// Warehouse
// -----------------------------------------------------------------------------

func (r *WMSRepository) CreateWarehouse(
	ctx context.Context,
	warehouse *model.Warehouse,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(warehouse).
		Error
}

func (r *WMSRepository) GetWarehouseByID(
	ctx context.Context,
	id uint,
) (*model.Warehouse, error) {
	var warehouse model.Warehouse

	err := r.db.Get().
		WithContext(ctx).
		Where("id = ?", id).
		First(&warehouse).
		Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, wms.ErrWarehouseNotFound
		// }
		return nil, err
	}

	return &warehouse, nil
}

// -----------------------------------------------------------------------------
// Warehouse Location
// -----------------------------------------------------------------------------

func (r *WMSRepository) CreateLocation(
	ctx context.Context,
	location *model.WarehouseLocation,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(location).
		Error
}

func (r *WMSRepository) GetLocationByID(
	ctx context.Context,
	id uint,
) (*model.WarehouseLocation, error) {
	var location model.WarehouseLocation

	err := r.db.Get().
		WithContext(ctx).
		Where("id = ?", id).
		First(&location).
		Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, wms.ErrLocationNotFound
		// }

		return nil, err
	}

	return &location, nil
}

// -----------------------------------------------------------------------------
// Device Inventory
// -----------------------------------------------------------------------------

func (r *WMSRepository) GetInventoryByDeviceID(
	ctx context.Context,
	deviceID uint,
) (*model.DeviceInventory, error) {
	var inventory model.DeviceInventory

	err := r.db.Get().
		WithContext(ctx).
		Where("device_id = ?", deviceID).
		First(&inventory).
		Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, wms.ErrInventoryNotFound
		// }

		return nil, err
	}

	return &inventory, nil
}

func (r *WMSRepository) GetInventoryByDeviceIDForUpdateTx(
	ctx context.Context,
	deviceID uint,
) (*model.DeviceInventory, error) {
	var inventory model.DeviceInventory

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("device_id = ?", deviceID).
		First(&inventory).
		Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, wms.ErrInventoryNotFound
		// }

		return nil, err
	}

	return &inventory, nil
}

func (r *WMSRepository) CreateInventoryTx(
	ctx context.Context,
	inventory *model.DeviceInventory,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(inventory).
		Error
}

func (r *WMSRepository) UpdateInventoryTx(
	ctx context.Context,
	inventory *model.DeviceInventory,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Save(inventory).
		Error
}

// -----------------------------------------------------------------------------
// Shipment
// -----------------------------------------------------------------------------

func (r *WMSRepository) CreateShipmentTx(
	ctx context.Context,
	shipment *model.Shipment,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(shipment).
		Error
}

func (r *WMSRepository) CreateShipmentItemsTx(
	ctx context.Context,
	items []*model.ShipmentItem,
) error {
	if len(items) == 0 {
		return nil
	}

	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		CreateInBatches(items, len(items)).
		Error
}

func (r *WMSRepository) GetShipmentByID(
	ctx context.Context,
	id uint,
) (*model.Shipment, error) {
	var shipment model.Shipment

	err := r.db.Get().
		WithContext(ctx).
		Where("id = ?", id).
		First(&shipment).
		Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, wms.ErrShipmentNotFound
		// }

		return nil, err
	}

	return &shipment, nil
}

func (r *WMSRepository) GetShipmentByIDForUpdateTx(
	ctx context.Context,
	id uint,
) (*model.Shipment, error) {
	var shipment model.Shipment

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("id = ?", id).
		First(&shipment).
		Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, wms.ErrShipmentNotFound
		// }

		return nil, err
	}

	return &shipment, nil
}

func (r *WMSRepository) ListShipmentItems(
	ctx context.Context,
	shipmentID uint,
) ([]*model.ShipmentItem, error) {
	var items []*model.ShipmentItem

	err := r.db.Get().
		WithContext(ctx).
		Where("shipment_id = ?", shipmentID).
		Order("id ASC").
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *WMSRepository) UpdateShipmentTx(
	ctx context.Context,
	shipment *model.Shipment,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Save(shipment).
		Error
}

// -----------------------------------------------------------------------------
// Tracking Events
// -----------------------------------------------------------------------------

func (r *WMSRepository) GetTrackingEventByExternalID(
	ctx context.Context,
	shipmentID uint,
	externalEventID string,
) (*model.ShipmentTrackingEvent, error) {
	var event model.ShipmentTrackingEvent

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"shipment_id = ? AND external_event_id = ?",
			shipmentID,
			externalEventID,
		).
		First(&event).
		Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, wms.ErrTrackingEventNotFound
		// }

		return nil, err
	}

	return &event, nil
}

func (r *WMSRepository) CreateTrackingEventTx(
	ctx context.Context,
	event *model.ShipmentTrackingEvent,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(event).
		Error
}

func (r *WMSRepository) ListTrackingEvents(
	ctx context.Context,
	shipmentID uint,
) ([]*model.ShipmentTrackingEvent, error) {
	var events []*model.ShipmentTrackingEvent

	err := r.db.Get().
		WithContext(ctx).
		Where("shipment_id = ?", shipmentID).
		Order("occurred_at ASC").
		Order("id ASC").
		Find(&events).
		Error

	if err != nil {
		return nil, err
	}

	return events, nil
}
