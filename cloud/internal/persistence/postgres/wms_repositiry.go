package postgres

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/wms"
	"github.com/boqrs/nexus/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type wmsRepository struct {
	db *database.DBProvider
}

func NewWMSRepository(db *database.DBProvider) wms.Repository {
	return &wmsRepository{
		db: db,
	}
}

// ============================================================
// Warehouse
// ============================================================

func (r *wmsRepository) CreateWarehouse(
	ctx context.Context,
	warehouse *model.Warehouse,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(warehouse).
		Error
}

func (r *wmsRepository) GetWarehouseByID(
	ctx context.Context,
	id uint,
) (*model.Warehouse, error) {
	var warehouse model.Warehouse

	err := r.db.Get().
		WithContext(ctx).
		First(&warehouse, id).
		Error
	if err != nil {
		return nil, err
	}

	return &warehouse, nil
}

// ============================================================
// Location
// ============================================================

func (r *wmsRepository) CreateLocation(
	ctx context.Context,
	location *model.WarehouseLocation,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(location).
		Error
}

func (r *wmsRepository) GetLocationByID(
	ctx context.Context,
	id uint,
) (*model.WarehouseLocation, error) {
	var location model.WarehouseLocation

	err := r.db.Get().
		WithContext(ctx).
		First(&location, id).
		Error
	if err != nil {
		return nil, err
	}

	return &location, nil
}

// ============================================================
// Inventory
// ============================================================

func (r *wmsRepository) GetInventoryByDeviceID(
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
		return nil, err
	}

	return &inventory, nil
}

func (r *wmsRepository) GetInventoryByDeviceIDForUpdateTx(
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
		return nil, err
	}

	return &inventory, nil
}

func (r *wmsRepository) CreateInventoryTx(
	ctx context.Context,
	inventory *model.DeviceInventory,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(inventory).
		Error
}

func (r *wmsRepository) UpdateInventoryTx(
	ctx context.Context,
	inventory *model.DeviceInventory,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Save(inventory).
		Error
}

// ============================================================
// Shipment
// ============================================================

func (r *wmsRepository) CreateShipmentTx(
	ctx context.Context,
	shipment *model.Shipment,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(shipment).
		Error
}

func (r *wmsRepository) CreateShipmentItemsTx(
	ctx context.Context,
	items []*model.ShipmentItem,
) error {
	if len(items) == 0 {
		return nil
	}

	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(&items).
		Error
}

func (r *wmsRepository) GetShipmentByID(
	ctx context.Context,
	id uint,
) (*model.Shipment, error) {
	var shipment model.Shipment

	err := r.db.Get().
		WithContext(ctx).
		First(&shipment, id).
		Error
	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *wmsRepository) GetShipmentByIDForUpdateTx(
	ctx context.Context,
	id uint,
) (*model.Shipment, error) {
	var shipment model.Shipment

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		First(&shipment, id).
		Error
	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *wmsRepository) ListShipmentItems(
	ctx context.Context,
	shipmentID uint,
) ([]*model.ShipmentItem, error) {
	var items []*model.ShipmentItem

	err := dbFromContext(ctx, r.db.Get()).
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

func (r *wmsRepository) UpdateShipmentTx(
	ctx context.Context,
	shipment *model.Shipment,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Save(shipment).
		Error
}

// ============================================================
// Tracking
// ============================================================

func (r *wmsRepository) GetTrackingEventByExternalID(
	ctx context.Context,
	shipmentID uint,
	externalEventID string,
) (*model.ShipmentTrackingEvent, error) {
	var event model.ShipmentTrackingEvent

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Where(
			"shipment_id = ? AND external_event_id = ?",
			shipmentID,
			externalEventID,
		).
		First(&event).
		Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *wmsRepository) CreateTrackingEventTx(
	ctx context.Context,
	event *model.ShipmentTrackingEvent,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(event).
		Error
}

func (r *wmsRepository) ListTrackingEvents(
	ctx context.Context,
	shipmentID uint,
) ([]*model.ShipmentTrackingEvent, error) {
	var events []*model.ShipmentTrackingEvent

	err := r.db.Get().
		WithContext(ctx).
		Where("shipment_id = ?", shipmentID).
		Order("occurred_at ASC, id ASC").
		Find(&events).
		Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

// Compile-time interface check.
var _ wms.Repository = (*wmsRepository)(nil)

// Keep gorm imported explicitly for repository-level error handling
// compatibility with the rest of the postgres package.
var _ = gorm.ErrRecordNotFound
