package wms

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	salesorderSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/salesorder"
)

var (
	ErrTenantNotFound = errors.New(
		"tenant id not found in context",
	)

	ErrWarehouseNotFound = errors.New(
		"warehouse not found",
	)

	ErrLocationNotFound = errors.New(
		"warehouse location not found",
	)

	ErrInventoryNotFound = errors.New(
		"device inventory not found",
	)

	ErrShipmentNotFound = errors.New(
		"shipment not found",
	)

	ErrTrackingEventNotFound = errors.New(
		"tracking event not found",
	)

	ErrDeviceNotFound = errors.New(
		"device not found",
	)

	ErrDeviceAlreadyInStock = errors.New(
		"device is already in stock",
	)

	ErrDeviceNotInStock = errors.New(
		"device is not in stock",
	)

	ErrDeviceReserved = errors.New(
		"device inventory is reserved",
	)

	ErrLocationNotBelongToWarehouse = errors.New(
		"location does not belong to warehouse",
	)

	ErrShipmentAlreadyShipped = errors.New(
		"shipment has already been shipped",
	)

	ErrShipmentCancelled = errors.New(
		"shipment has been cancelled",
	)

	ErrShipmentDelivered = errors.New(
		"shipment has already been delivered",
	)

	ErrInvalidShipmentStatus = errors.New(
		"invalid shipment status",
	)

	ErrDuplicateDevice = errors.New(
		"duplicate device in shipment",
	)

	ErrEmptyShipment = errors.New(
		"shipment must contain at least one device",
	)

	ErrInvalidTrackingEvent = errors.New(
		"invalid tracking event",
	)

	ErrSalesOrderNotFound = errors.New(
		"sales order not found",
	)

	ErrSalesOrderNotConfirmed = errors.New(
		"sales order is not confirmed",
	)

	ErrSalesOrderItemNotFound = errors.New(
		"sales order item not found",
	)

	ErrSalesOrderItemMismatch = errors.New(
		"sales order item does not belong to shipment sales order",
	)

	ErrSalesOrderItemRequired = errors.New(
		"sales order item is required for sales order shipment",
	)

	ErrSalesOrderFulfillmentExceeded = errors.New(
		"sales order fulfillment quantity exceeds ordered quantity",
	)

	ErrSalesOrderShipmentMismatch = errors.New(
		"shipment items do not belong to the shipment sales order",
	)
)

type service struct {
	uow         UnitOfWork
	repository  Repository
	salesOrders salesorderSrv.Service
}

func NewService(
	uow UnitOfWork,
	repository Repository,
	salesOrders salesorderSrv.Service,
) Service {
	return &service{
		uow:         uow,
		repository:  repository,
		salesOrders: salesOrders,
	}
}

// -----------------------------------------------------------------------------
// Warehouse
// -----------------------------------------------------------------------------

func (s *service) CreateWarehouse(
	ctx context.Context,
	req *CreateWarehouseRequest,
) (*WarehouseResponse, error) {
	if req == nil {
		return nil, ErrWarehouseNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrTenantNotFound
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)

	if code == "" || name == "" {
		return nil, ErrWarehouseNotFound
	}

	warehouse := &model.Warehouse{
		TenantID: tenantID,
		Code:     code,
		Name:     name,
		Address:  strings.TrimSpace(req.Address),
	}

	if err := s.repository.CreateWarehouse(
		ctx,
		warehouse,
	); err != nil {
		return nil, err
	}

	return warehouseToResponse(warehouse), nil
}

func (s *service) GetWarehouse(
	ctx context.Context,
	id uint,
) (*WarehouseResponse, error) {
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrTenantNotFound
	}

	warehouse, err := s.repository.GetWarehouseByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return nil, err
	}

	return warehouseToResponse(warehouse), nil
}

// -----------------------------------------------------------------------------
// Location
// -----------------------------------------------------------------------------

func (s *service) CreateLocation(
	ctx context.Context,
	req *CreateLocationRequest,
) (*LocationResponse, error) {
	if req == nil {
		return nil, ErrLocationNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrTenantNotFound
	}

	if req.WarehouseID == 0 {
		return nil, ErrWarehouseNotFound
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)

	if code == "" || name == "" {
		return nil, ErrLocationNotFound
	}

	if _, err := s.repository.GetWarehouseByID(
		ctx,
		tenantID,
		req.WarehouseID,
	); err != nil {
		return nil, err
	}

	location := &model.WarehouseLocation{
		WarehouseID: req.WarehouseID,
		Code:        code,
		Name:        name,
	}

	if err := s.repository.CreateLocation(
		ctx,
		location,
	); err != nil {
		return nil, err
	}

	return locationToResponse(location), nil
}

// -----------------------------------------------------------------------------
// Inventory
// -----------------------------------------------------------------------------

func (s *service) GetDeviceInventory(
	ctx context.Context,
	deviceID uint,
) (*InventoryResponse, error) {
	if deviceID == 0 {
		return nil, ErrDeviceNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrTenantNotFound
	}

	inventory, err := s.repository.GetInventoryByDeviceID(
		ctx,
		tenantID,
		deviceID,
	)
	if err != nil {
		return nil, err
	}

	return inventoryToResponse(inventory), nil
}

func (s *service) StockIn(
	ctx context.Context,
	req *StockInRequest,
) (*InventoryResponse, error) {
	if req == nil || req.DeviceID == 0 {
		return nil, ErrDeviceNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrTenantNotFound
	}

	if req.WarehouseID == 0 {
		return nil, ErrWarehouseNotFound
	}

	if req.LocationID == 0 {
		return nil, ErrLocationNotFound
	}

	deviceBelongsToTenant, err :=
		s.repository.DeviceBelongsToTenant(
			ctx,
			tenantID,
			req.DeviceID,
		)
	if err != nil {
		return nil, err
	}

	if !deviceBelongsToTenant {
		return nil, ErrDeviceNotFound
	}

	if _, err := s.repository.GetWarehouseByID(
		ctx,
		tenantID,
		req.WarehouseID,
	); err != nil {
		return nil, err
	}

	location, err := s.repository.GetLocationByID(
		ctx,
		tenantID,
		req.LocationID,
	)
	if err != nil {
		return nil, err
	}

	if location.WarehouseID != req.WarehouseID {
		return nil, ErrLocationNotBelongToWarehouse
	}

	inboundAt := time.Now().UTC()

	if req.InboundAt != nil {
		inboundAt = req.InboundAt.UTC()
	}

	var result *model.DeviceInventory

	err = s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			inventory, err :=
				s.repository.GetInventoryByDeviceIDForUpdateTx(
					txCtx,
					tenantID,
					req.DeviceID,
				)

			if err != nil {
				if !errors.Is(err, ErrInventoryNotFound) {
					return err
				}

				inventory = &model.DeviceInventory{
					DeviceID:    req.DeviceID,
					WarehouseID: req.WarehouseID,
					LocationID:  req.LocationID,
					Status:      model.InventoryStatusInStock,
					InboundAt:   inboundAt,
				}

				if err := s.repository.CreateInventoryTx(
					txCtx,
					inventory,
				); err != nil {
					return err
				}

				result = inventory
				return nil
			}

			switch inventory.Status {
			case model.InventoryStatusInStock:
				return ErrDeviceAlreadyInStock

			case model.InventoryStatusReserved:
				return ErrDeviceReserved

			case model.InventoryStatusShipped:
				inventory.WarehouseID = req.WarehouseID
				inventory.LocationID = req.LocationID
				inventory.Status = model.InventoryStatusInStock
				inventory.InboundAt = inboundAt
				inventory.OutboundAt = nil

				if err := s.repository.UpdateInventoryTx(
					txCtx,
					inventory,
				); err != nil {
					return err
				}

				result = inventory
				return nil

			default:
				return ErrDeviceNotInStock
			}
		},
	)

	if err != nil {
		return nil, err
	}

	return inventoryToResponse(result), nil
}

// -----------------------------------------------------------------------------
// Shipment
// -----------------------------------------------------------------------------

func (s *service) CreateShipment(
	ctx context.Context,
	req *CreateShipmentRequest,
) (*ShipmentResponse, error) {
	if req == nil {
		return nil, ErrEmptyShipment
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrTenantNotFound
	}

	carrier := strings.TrimSpace(req.Carrier)
	trackingNumber := strings.TrimSpace(req.TrackingNumber)

	if carrier == "" || trackingNumber == "" {
		return nil, ErrInvalidShipmentStatus
	}

	items := normalizeShipmentItems(req)

	if len(items) == 0 {
		return nil, ErrEmptyShipment
	}

	seen := make(map[uint]struct{}, len(items))
	deviceIDs := make([]uint, 0, len(items))

	for _, item := range items {
		if item.DeviceID == 0 {
			return nil, ErrDeviceNotFound
		}

		if _, exists := seen[item.DeviceID]; exists {
			return nil, ErrDuplicateDevice
		}

		seen[item.DeviceID] = struct{}{}
		deviceIDs = append(deviceIDs, item.DeviceID)
	}

	// A shipment without a sales order must not contain
	// sales-order item references.
	if req.SalesOrderID == nil {
		for _, item := range items {
			if item.SalesOrderItemID != nil {
				return nil, ErrSalesOrderShipmentMismatch
			}
		}
	} else {
		// A shipment linked to a sales order must identify
		// the corresponding sales-order item for every device.
		for _, item := range items {
			if item.SalesOrderItemID == nil ||
				*item.SalesOrderItemID == 0 {
				return nil, ErrSalesOrderItemRequired
			}
		}
	}

	sort.Slice(
		deviceIDs,
		func(i, j int) bool {
			return deviceIDs[i] < deviceIDs[j]
		},
	)

	shipment := &model.Shipment{
		TenantID: tenantID,

		SalesOrderID: req.SalesOrderID,

		ExternalOrderID: strings.TrimSpace(
			req.ExternalOrderID,
		),

		Carrier:        carrier,
		TrackingNumber: trackingNumber,
		Status:         model.ShipmentStatusCreated,
	}

	shipmentItems := make(
		[]*model.ShipmentItem,
		0,
		len(items),
	)

	err := s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			// Lock and reserve the corresponding sales-order
			// quantities through the SalesOrder business service.
			//
			// The SalesOrder service must execute its repository
			// operations against the transaction carried by txCtx.
			if req.SalesOrderID != nil {
				if s.salesOrders == nil {
					return ErrSalesOrderNotFound
				}

				counts := shipmentSalesOrderItemCounts(items)

				if err := s.salesOrders.ReserveForShipment(
					txCtx,
					*req.SalesOrderID,
					counts,
				); err != nil {
					return mapSalesOrderError(err)
				}
			}

			lockedInventories := make(
				[]*model.DeviceInventory,
				0,
				len(deviceIDs),
			)

			// Always lock inventories in DeviceID order.
			// This keeps the locking order deterministic and
			// avoids unnecessary deadlock risk.
			for _, deviceID := range deviceIDs {
				inventory, err :=
					s.repository.GetInventoryByDeviceIDForUpdateTx(
						txCtx,
						tenantID,
						deviceID,
					)

				if err != nil {
					if errors.Is(err, ErrInventoryNotFound) {
						return ErrDeviceNotInStock
					}

					return err
				}

				if inventory.Status != model.InventoryStatusInStock {
					return ErrDeviceNotInStock
				}

				lockedInventories = append(
					lockedInventories,
					inventory,
				)
			}

			if err := s.repository.CreateShipmentTx(
				txCtx,
				shipment,
			); err != nil {
				return err
			}

			for _, item := range items {
				shipmentItems = append(
					shipmentItems,
					&model.ShipmentItem{
						ShipmentID:       shipment.ID,
						DeviceID:         item.DeviceID,
						SalesOrderItemID: item.SalesOrderItemID,
					},
				)
			}

			if err := s.repository.CreateShipmentItemsTx(
				txCtx,
				shipmentItems,
			); err != nil {
				return err
			}

			for _, inventory := range lockedInventories {
				inventory.Status =
					model.InventoryStatusReserved

				if err := s.repository.UpdateInventoryTx(
					txCtx,
					inventory,
				); err != nil {
					return err
				}
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &ShipmentResponse{
		ID:              shipment.ID,
		SalesOrderID:    shipment.SalesOrderID,
		ExternalOrderID: shipment.ExternalOrderID,
		Carrier:         shipment.Carrier,
		TrackingNumber:  shipment.TrackingNumber,
		Status:          shipment.Status.String(),
		Items:           shipmentItemsToResponse(shipmentItems),
		CreatedAt:       shipment.CreatedAt,
		UpdatedAt:       shipment.UpdatedAt,
	}, nil
}

// -----------------------------------------------------------------------------
// Stock Out
// -----------------------------------------------------------------------------

func (s *service) StockOut(
	ctx context.Context,
	shipmentID uint,
) error {
	if shipmentID == 0 {
		return ErrShipmentNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return ErrTenantNotFound
	}

	return s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			shipment, err :=
				s.repository.GetShipmentByIDForUpdateTx(
					txCtx,
					tenantID,
					shipmentID,
				)

			if err != nil {
				return err
			}

			switch shipment.Status {
			case model.ShipmentStatusCreated:

			case model.ShipmentStatusCancelled:
				return ErrShipmentCancelled

			case model.ShipmentStatusDelivered:
				return ErrShipmentDelivered

			default:
				return ErrShipmentAlreadyShipped
			}

			items, err := s.repository.ListShipmentItems(
				txCtx,
				tenantID,
				shipmentID,
			)
			if err != nil {
				return err
			}

			if len(items) == 0 {
				return ErrEmptyShipment
			}

			sort.Slice(
				items,
				func(i, j int) bool {
					return items[i].DeviceID < items[j].DeviceID
				},
			)

			now := time.Now().UTC()

			for _, item := range items {
				inventory, err :=
					s.repository.GetInventoryByDeviceIDForUpdateTx(
						txCtx,
						tenantID,
						item.DeviceID,
					)

				if err != nil {
					if errors.Is(err, ErrInventoryNotFound) {
						return ErrDeviceNotInStock
					}

					return err
				}

				if inventory.Status != model.InventoryStatusReserved {
					return ErrDeviceNotInStock
				}

				inventory.Status =
					model.InventoryStatusShipped

				inventory.OutboundAt = &now

				if err := s.repository.UpdateInventoryTx(
					txCtx,
					inventory,
				); err != nil {
					return err
				}
			}

			shipment.Status =
				model.ShipmentStatusInTransit

			shipment.ShippedAt = &now

			return s.repository.UpdateShipmentTx(
				txCtx,
				shipment,
			)
		},
	)
}

// -----------------------------------------------------------------------------
// Shipment Query
// -----------------------------------------------------------------------------

func (s *service) GetShipment(
	ctx context.Context,
	shipmentID uint,
) (*ShipmentResponse, error) {
	if shipmentID == 0 {
		return nil, ErrShipmentNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrTenantNotFound
	}

	shipment, err := s.repository.GetShipmentByID(
		ctx,
		tenantID,
		shipmentID,
	)
	if err != nil {
		return nil, err
	}

	items, err := s.repository.ListShipmentItems(
		ctx,
		tenantID,
		shipmentID,
	)
	if err != nil {
		return nil, err
	}

	return &ShipmentResponse{
		ID:              shipment.ID,
		SalesOrderID:    shipment.SalesOrderID,
		ExternalOrderID: shipment.ExternalOrderID,
		Carrier:         shipment.Carrier,
		TrackingNumber:  shipment.TrackingNumber,
		Status:          shipment.Status.String(),
		ShippedAt:       shipment.ShippedAt,
		DeliveredAt:     shipment.DeliveredAt,
		Items:           shipmentItemsToResponse(items),
		CreatedAt:       shipment.CreatedAt,
		UpdatedAt:       shipment.UpdatedAt,
	}, nil
}

// -----------------------------------------------------------------------------
// Tracking
// -----------------------------------------------------------------------------

func (s *service) AddTrackingEvent(
	ctx context.Context,
	shipmentID uint,
	req *TrackingEventRequest,
) error {
	if shipmentID == 0 {
		return ErrShipmentNotFound
	}

	if req == nil {
		return ErrInvalidTrackingEvent
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return ErrTenantNotFound
	}

	status, ok := req.Status.ToModel()
	if !ok {
		return ErrInvalidShipmentStatus
	}

	externalEventID := strings.TrimSpace(
		req.ExternalEventID,
	)

	if externalEventID == "" {
		return ErrInvalidTrackingEvent
	}

	occurredAt := time.Now().UTC()

	if req.OccurredAt != nil {
		occurredAt = req.OccurredAt.UTC()
	}

	return s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			shipment, err :=
				s.repository.GetShipmentByIDForUpdateTx(
					txCtx,
					tenantID,
					shipmentID,
				)
			if err != nil {
				return err
			}

			existing, err :=
				s.repository.GetTrackingEventByExternalID(
					txCtx,
					tenantID,
					shipmentID,
					externalEventID,
				)

			if err == nil {
				_ = existing
				return nil
			}

			if !errors.Is(err, ErrTrackingEventNotFound) {
				return err
			}

			nextStatus := model.ShipmentStatus(status)

			if !isValidShipmentStatusTransition(
				shipment.Status,
				nextStatus,
			) {
				return ErrInvalidShipmentStatus
			}

			event := &model.ShipmentTrackingEvent{
				ShipmentID:      shipmentID,
				ExternalEventID: externalEventID,
				Status:          nextStatus,
				OccurredAt:      occurredAt,
				Location:        strings.TrimSpace(req.Location),
				Description:     strings.TrimSpace(req.Description),
				CreatedAt:       time.Now().UTC(),
			}

			if err := s.repository.CreateTrackingEventTx(
				txCtx,
				event,
			); err != nil {
				return err
			}

			previousStatus := shipment.Status

			shipment.Status = nextStatus
			shipment.UpdatedAt = time.Now().UTC()

			if nextStatus == model.ShipmentStatusDelivered {
				deliveredAt := occurredAt
				shipment.DeliveredAt = &deliveredAt

				if err := s.handleShipmentDeliveredTx(
					txCtx,
					shipment,
					previousStatus,
				); err != nil {
					return err
				}
			}

			return s.repository.UpdateShipmentTx(
				txCtx,
				shipment,
			)
		},
	)
}

func (s *service) ListTrackingEvents(
	ctx context.Context,
	shipmentID uint,
) ([]*TrackingEventResponse, error) {
	if shipmentID == 0 {
		return nil, ErrShipmentNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrTenantNotFound
	}

	if _, err := s.repository.GetShipmentByID(
		ctx,
		tenantID,
		shipmentID,
	); err != nil {
		return nil, err
	}

	events, err := s.repository.ListTrackingEvents(
		ctx,
		tenantID,
		shipmentID,
	)
	if err != nil {
		return nil, err
	}

	result := make(
		[]*TrackingEventResponse,
		0,
		len(events),
	)

	for _, event := range events {
		result = append(
			result,
			trackingEventToResponse(event),
		)
	}

	return result, nil
}

// -----------------------------------------------------------------------------
// Shipment Cancellation
// -----------------------------------------------------------------------------

func (s *service) CancelShipment(
	ctx context.Context,
	shipmentID uint,
) error {
	if shipmentID == 0 {
		return ErrShipmentNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return ErrTenantNotFound
	}

	return s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			shipment, err :=
				s.repository.GetShipmentByIDForUpdateTx(
					txCtx,
					tenantID,
					shipmentID,
				)

			if err != nil {
				return err
			}

			switch shipment.Status {
			case model.ShipmentStatusCreated:

			case model.ShipmentStatusCancelled:
				return nil

			case model.ShipmentStatusInTransit,
				model.ShipmentStatusOutForDelivery,
				model.ShipmentStatusDelivered,
				model.ShipmentStatusException:
				return ErrShipmentAlreadyShipped

			default:
				return ErrInvalidShipmentStatus
			}

			items, err := s.repository.ListShipmentItems(
				txCtx,
				tenantID,
				shipmentID,
			)
			if err != nil {
				return err
			}

			sort.Slice(
				items,
				func(i, j int) bool {
					return items[i].DeviceID < items[j].DeviceID
				},
			)

			now := time.Now().UTC()

			for _, item := range items {
				inventory, err :=
					s.repository.GetInventoryByDeviceIDForUpdateTx(
						txCtx,
						tenantID,
						item.DeviceID,
					)

				if err != nil {
					return err
				}

				if inventory.Status != model.InventoryStatusReserved {
					return ErrDeviceNotInStock
				}

				inventory.Status =
					model.InventoryStatusInStock

				inventory.OutboundAt = nil
				inventory.UpdatedAt = now

				if err := s.repository.UpdateInventoryTx(
					txCtx,
					inventory,
				); err != nil {
					return err
				}
			}

			if shipment.SalesOrderID != nil {
				if s.salesOrders == nil {
					return ErrSalesOrderNotFound
				}

				counts := shipmentSalesOrderItemCountsFromShipmentItems(
					items,
				)

				if err := s.salesOrders.ReleaseShipmentReservation(
					txCtx,
					*shipment.SalesOrderID,
					counts,
				); err != nil {
					return mapSalesOrderError(err)
				}
			}

			shipment.Status =
				model.ShipmentStatusCancelled

			shipment.UpdatedAt = now

			return s.repository.UpdateShipmentTx(
				txCtx,
				shipment,
			)
		},
	)
}

// -----------------------------------------------------------------------------
// Sales Order Fulfillment
// -----------------------------------------------------------------------------

func (s *service) handleShipmentDeliveredTx(
	ctx context.Context,
	shipment *model.Shipment,
	previousStatus model.ShipmentStatus,
) error {
	if shipment.SalesOrderID == nil {
		return nil
	}

	if previousStatus != model.ShipmentStatusInTransit &&
		previousStatus != model.ShipmentStatusOutForDelivery &&
		previousStatus != model.ShipmentStatusException {
		return ErrInvalidShipmentStatus
	}

	if s.salesOrders == nil {
		return ErrSalesOrderNotFound
	}

	items, err := s.repository.ListShipmentItems(
		ctx,
		shipment.TenantID,
		shipment.ID,
	)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return ErrEmptyShipment
	}

	counts := shipmentSalesOrderItemCountsFromShipmentItems(
		items,
	)

	if len(counts) == 0 {
		return ErrSalesOrderItemRequired
	}

	if err := s.salesOrders.FulfillShipment(
		ctx,
		*shipment.SalesOrderID,
		counts,
	); err != nil {
		return mapSalesOrderError(err)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func normalizeShipmentItems(
	req *CreateShipmentRequest,
) []CreateShipmentItemRequest {
	if len(req.Items) > 0 {
		return append(
			[]CreateShipmentItemRequest(nil),
			req.Items...,
		)
	}

	if len(req.DeviceIDs) == 0 {
		return nil
	}

	items := make(
		[]CreateShipmentItemRequest,
		0,
		len(req.DeviceIDs),
	)

	for _, deviceID := range req.DeviceIDs {
		items = append(
			items,
			CreateShipmentItemRequest{
				DeviceID: deviceID,
			},
		)
	}

	return items
}

func shipmentSalesOrderItemCounts(
	items []CreateShipmentItemRequest,
) map[uint]int64 {
	result := make(map[uint]int64)

	for _, item := range items {
		if item.SalesOrderItemID == nil ||
			*item.SalesOrderItemID == 0 {
			continue
		}

		result[*item.SalesOrderItemID]++
	}

	return result
}

func shipmentSalesOrderItemCountsFromShipmentItems(
	items []*model.ShipmentItem,
) map[uint]int64 {
	result := make(map[uint]int64)

	for _, item := range items {
		if item == nil ||
			item.SalesOrderItemID == nil ||
			*item.SalesOrderItemID == 0 {
			continue
		}

		result[*item.SalesOrderItemID]++
	}

	return result
}

func mapSalesOrderError(err error) error {
	if errors.Is(err, salesorderSrv.ErrSalesOrderNotFound) {
		return ErrSalesOrderNotFound
	}

	return err
}

func warehouseToResponse(
	warehouse *model.Warehouse,
) *WarehouseResponse {
	return &WarehouseResponse{
		ID:        warehouse.ID,
		Code:      warehouse.Code,
		Name:      warehouse.Name,
		Address:   warehouse.Address,
		CreatedAt: warehouse.CreatedAt,
		UpdatedAt: warehouse.UpdatedAt,
	}
}

func locationToResponse(
	location *model.WarehouseLocation,
) *LocationResponse {
	return &LocationResponse{
		ID:          location.ID,
		WarehouseID: location.WarehouseID,
		Code:        location.Code,
		Name:        location.Name,
		CreatedAt:   location.CreatedAt,
		UpdatedAt:   location.UpdatedAt,
	}
}

func inventoryToResponse(
	inventory *model.DeviceInventory,
) *InventoryResponse {
	return &InventoryResponse{
		ID:          inventory.ID,
		DeviceID:    inventory.DeviceID,
		WarehouseID: inventory.WarehouseID,
		LocationID:  inventory.LocationID,
		Status:      inventory.Status.String(),
		InboundAt:   inventory.InboundAt,
		OutboundAt:  inventory.OutboundAt,
		CreatedAt:   inventory.CreatedAt,
		UpdatedAt:   inventory.UpdatedAt,
	}
}

func shipmentItemsToResponse(
	items []*model.ShipmentItem,
) []*ShipmentItemResponse {
	result := make(
		[]*ShipmentItemResponse,
		0,
		len(items),
	)

	for _, item := range items {
		if item == nil {
			continue
		}

		result = append(
			result,
			&ShipmentItemResponse{
				ID:               item.ID,
				ShipmentID:       item.ShipmentID,
				DeviceID:         item.DeviceID,
				SalesOrderItemID: item.SalesOrderItemID,
				CreatedAt:        item.CreatedAt,
			},
		)
	}

	return result
}

func trackingEventToResponse(
	event *model.ShipmentTrackingEvent,
) *TrackingEventResponse {
	return &TrackingEventResponse{
		ID:              event.ID,
		ShipmentID:      event.ShipmentID,
		ExternalEventID: event.ExternalEventID,
		Status:          event.Status.String(),
		OccurredAt:      event.OccurredAt,
		Location:        event.Location,
		Description:     event.Description,
		CreatedAt:       event.CreatedAt,
	}
}

func isValidShipmentStatusTransition(
	current model.ShipmentStatus,
	next model.ShipmentStatus,
) bool {
	switch current {
	case model.ShipmentStatusCreated:
		return next == model.ShipmentStatusInTransit

	case model.ShipmentStatusInTransit:
		return next == model.ShipmentStatusOutForDelivery ||
			next == model.ShipmentStatusException

	case model.ShipmentStatusOutForDelivery:
		return next == model.ShipmentStatusDelivered ||
			next == model.ShipmentStatusException

	case model.ShipmentStatusException:
		return next == model.ShipmentStatusInTransit ||
			next == model.ShipmentStatusOutForDelivery ||
			next == model.ShipmentStatusDelivered

	case model.ShipmentStatusDelivered,
		model.ShipmentStatusCancelled:
		return false

	default:
		return false
	}
}
