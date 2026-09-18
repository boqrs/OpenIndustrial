package postgres

import (
	"context"
	"errors"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/wms"
	"github.com/boqrs/nexus/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type wmsRepository struct {
	db *database.DBProvider
}

func NewWMSRepository(
	db *database.DBProvider,
) wms.Repository {
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
	tenantID uuid.UUID,
	id uint,
) (*model.Warehouse, error) {
	var warehouse model.Warehouse

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND id = ?",
			tenantID,
			id,
		).
		First(&warehouse).
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
	tenantID uuid.UUID,
	id uint,
) (*model.WarehouseLocation, error) {
	var location model.WarehouseLocation

	err := r.db.Get().
		WithContext(ctx).
		Table("warehouse_locations").
		Joins(
			"JOIN warehouses ON warehouses.id = warehouse_locations.warehouse_id",
		).
		Where(
			"warehouse_locations.id = ? AND warehouses.tenant_id = ?",
			id,
			tenantID,
		).
		First(&location).
		Error

	if err != nil {
		return nil, err
	}

	return &location, nil
}

// ============================================================
// Inventory
// ============================================================

// GetInventoryByDeviceID returns the current inventory record
// for a device.
//
// Device inventory itself does not carry TenantID. Tenant isolation
// is therefore performed through:
//
//	device_inventories -> devices -> resources
func (r *wmsRepository) GetInventoryByDeviceID(
	ctx context.Context,
	tenantID uuid.UUID,
	deviceID uint,
) (*model.DeviceInventory, error) {
	var inventory model.DeviceInventory

	err := r.db.Get().
		WithContext(ctx).
		Table("device_inventories").
		Joins(
			"JOIN devices ON devices.id = device_inventories.device_id",
		).
		Joins(
			"JOIN resources ON resources.id = devices.resource_id",
		).
		Where(
			"device_inventories.device_id = ? AND resources.tenant_id = ?",
			deviceID,
			tenantID,
		).
		First(&inventory).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, wms.ErrInventoryNotFound
		}

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

// DeviceBelongsToTenant verifies that a Device belongs to the
// authenticated tenant.
//
// Device does not contain TenantID directly.
// Tenant ownership is inherited through:
//
//	Device -> Resource -> Tenant
//
// This method intentionally lives in WMS repository rather than
// calling device.Repository.GetByID(), because tenant isolation
// is a persistence concern and must be enforced by the query itself.
func (r *wmsRepository) DeviceBelongsToTenant(
	ctx context.Context,
	tenantID uuid.UUID,
	deviceID uint,
) (bool, error) {
	var count int64

	err := r.db.Get().
		WithContext(ctx).
		Table("devices").
		Joins(
			"JOIN resources ON resources.id = devices.resource_id",
		).
		Where(
			"devices.id = ? AND resources.tenant_id = ?",
			deviceID,
			tenantID,
		).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
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
	tenantID uuid.UUID,
	id uint,
) (*model.Shipment, error) {
	var shipment model.Shipment

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND id = ?",
			tenantID,
			id,
		).
		First(&shipment).
		Error

	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *wmsRepository) GetShipmentByIDForUpdateTx(
	ctx context.Context,
	tenantID uuid.UUID,
	id uint,
) (*model.Shipment, error) {
	var shipment model.Shipment

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Where(
			"tenant_id = ? AND id = ?",
			tenantID,
			id,
		).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		First(&shipment).
		Error

	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *wmsRepository) ListShipmentItems(
	ctx context.Context,
	tenantID uuid.UUID,
	shipmentID uint,
) ([]*model.ShipmentItem, error) {
	var items []*model.ShipmentItem

	err := r.db.Get().
		WithContext(ctx).
		Table("shipment_items").
		Joins(
			"JOIN shipments ON shipments.id = shipment_items.shipment_id",
		).
		Where(
			"shipment_items.shipment_id = ? AND shipments.tenant_id = ?",
			shipmentID,
			tenantID,
		).
		Order("shipment_items.id ASC").
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
	tenantID uuid.UUID,
	shipmentID uint,
	externalEventID string,
) (*model.ShipmentTrackingEvent, error) {
	var event model.ShipmentTrackingEvent

	err := r.db.Get().
		WithContext(ctx).
		Table("shipment_tracking_events").
		Joins(
			"JOIN shipments ON shipments.id = shipment_tracking_events.shipment_id",
		).
		Where(
			"shipment_tracking_events.shipment_id = ? "+
				"AND shipment_tracking_events.external_event_id = ? "+
				"AND shipments.tenant_id = ?",
			shipmentID,
			externalEventID,
			tenantID,
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
	tenantID uuid.UUID,
	shipmentID uint,
) ([]*model.ShipmentTrackingEvent, error) {
	var events []*model.ShipmentTrackingEvent

	err := r.db.Get().
		WithContext(ctx).
		Table("shipment_tracking_events").
		Joins(
			"JOIN shipments ON shipments.id = shipment_tracking_events.shipment_id",
		).
		Where(
			"shipment_tracking_events.shipment_id = ? "+
				"AND shipments.tenant_id = ?",
			shipmentID,
			tenantID,
		).
		Order(
			"shipment_tracking_events.occurred_at ASC, " +
				"shipment_tracking_events.id ASC",
		).
		Find(&events).
		Error

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *wmsRepository) GetInventoryByDeviceIDForUpdateTx(
	ctx context.Context,
	tenantID uuid.UUID,
	deviceID uint,
) (*model.DeviceInventory, error) {
	var inventory model.DeviceInventory

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Table("device_inventories").
		Joins(
			"JOIN devices ON devices.id = device_inventories.device_id",
		).
		Joins(
			"JOIN resources ON resources.id = devices.resource_id",
		).
		Where(
			"device_inventories.device_id = ? AND resources.tenant_id = ?",
			deviceID,
			tenantID,
		).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		First(&inventory).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, wms.ErrInventoryNotFound
		}

		return nil, err
	}

	return &inventory, nil
}

// ============================================================
// Compile-time interface check
// ============================================================

var _ wms.Repository = (*wmsRepository)(nil)

// Keep the GORM dependency explicit for compatibility with
// the postgres package.
var _ = gorm.ErrRecordNotFound
