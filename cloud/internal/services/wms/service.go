package wms

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
)

var (
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
)

type service struct {
	uow        UnitOfWork
	repository Repository
	devices    device.Repository
}

func NewService(
	uow UnitOfWork,
	repository Repository,
	devices device.Repository,
) Service {
	return &service{
		uow:        uow,
		repository: repository,
		devices:    devices,
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

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)

	if code == "" || name == "" {
		return nil, ErrWarehouseNotFound
	}

	warehouse := &model.Warehouse{
		Code:    code,
		Name:    name,
		Address: strings.TrimSpace(req.Address),
	}

	if err := s.repository.CreateWarehouse(ctx, warehouse); err != nil {
		return nil, err
	}

	return warehouseToResponse(warehouse), nil
}

func (s *service) GetWarehouse(
	ctx context.Context,
	id uint,
) (*WarehouseResponse, error) {

	warehouse, err := s.repository.GetWarehouseByID(ctx, id)
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

	if req.WarehouseID == 0 {
		return nil, ErrWarehouseNotFound
	}

	if req.Code == "" || req.Name == "" {
		return nil, ErrLocationNotFound
	}

	if _, err := s.repository.GetWarehouseByID(
		ctx,
		req.WarehouseID,
	); err != nil {
		return nil, err
	}

	location := &model.WarehouseLocation{
		WarehouseID: req.WarehouseID,
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
	}

	if err := s.repository.CreateLocation(ctx, location); err != nil {
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

	inventory, err := s.repository.GetInventoryByDeviceID(
		ctx,
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

	if req.WarehouseID == 0 {
		return nil, ErrWarehouseNotFound
	}

	if req.LocationID == 0 {
		return nil, ErrLocationNotFound
	}

	// WMS never creates Device.
	// It only verifies that MES has already created it.
	if _, err := s.devices.GetByID(
		ctx,
		req.DeviceID,
	); err != nil {
		return nil, ErrDeviceNotFound
	}

	if _, err := s.repository.GetWarehouseByID(
		ctx,
		req.WarehouseID,
	); err != nil {
		return nil, err
	}

	location, err := s.repository.GetLocationByID(
		ctx,
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

			inventory, err := s.repository.
				GetInventoryByDeviceIDForUpdateTx(
					txCtx,
					req.DeviceID,
				)

			if err != nil {

				if err != ErrInventoryNotFound {
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

			if inventory.Status == model.InventoryStatusInStock {
				return ErrDeviceAlreadyInStock
			}

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

	if req == nil || len(req.DeviceIDs) == 0 {
		return nil, ErrEmptyShipment
	}

	carrier := strings.TrimSpace(req.Carrier)
	trackingNumber := strings.TrimSpace(req.TrackingNumber)

	if carrier == "" || trackingNumber == "" {
		return nil, ErrInvalidShipmentStatus
	}

	deviceIDs := make(
		[]uint,
		0,
		len(req.DeviceIDs),
	)

	seen := make(
		map[uint]struct{},
		len(req.DeviceIDs),
	)

	for _, deviceID := range req.DeviceIDs {

		if deviceID == 0 {
			return nil, ErrDeviceNotFound
		}

		if _, exists := seen[deviceID]; exists {
			return nil, ErrDuplicateDevice
		}

		seen[deviceID] = struct{}{}

		// Device must already exist.
		if _, err := s.devices.GetByID(
			ctx,
			deviceID,
		); err != nil {
			return nil, ErrDeviceNotFound
		}

		// A shipment is only allowed to be created
		// for devices that are currently in stock.
		inventory, err := s.repository.GetInventoryByDeviceID(
			ctx,
			deviceID,
		)
		if err != nil {
			return nil, err
		}

		if inventory.Status != model.InventoryStatusInStock {
			return nil, ErrDeviceNotInStock
		}

		deviceIDs = append(
			deviceIDs,
			deviceID,
		)
	}

	shipment := &model.Shipment{
		ExternalOrderID: strings.TrimSpace(
			req.ExternalOrderID,
		),
		Carrier:        carrier,
		TrackingNumber: trackingNumber,
		Status:         model.ShipmentStatusCreated,
	}

	items := make(
		[]*model.ShipmentItem,
		0,
		len(deviceIDs),
	)

	err := s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {

			if err := s.repository.CreateShipmentTx(
				txCtx,
				shipment,
			); err != nil {
				return err
			}

			for _, deviceID := range deviceIDs {

				items = append(
					items,
					&model.ShipmentItem{
						ShipmentID: shipment.ID,
						DeviceID:   deviceID,
					},
				)
			}

			return s.repository.CreateShipmentItemsTx(
				txCtx,
				items,
			)
		},
	)

	if err != nil {
		return nil, err
	}

	return &ShipmentResponse{
		ID:              shipment.ID,
		ExternalOrderID: shipment.ExternalOrderID,
		Carrier:         shipment.Carrier,
		TrackingNumber:  shipment.TrackingNumber,
		Status:          shipment.Status.String(),
		Items:           shipmentItemsToResponse(items),
		CreatedAt:       shipment.CreatedAt,
		UpdatedAt:       shipment.UpdatedAt,
	}, nil
}

// StockOut changes:
//
// Shipment:
//
//	CREATED -> IN_TRANSIT
//
// DeviceInventory:
//
//	IN_STOCK -> SHIPPED
//
// The entire operation is transactional.
func (s *service) StockOut(
	ctx context.Context,
	shipmentID uint,
) error {

	return s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {

			shipment, err := s.repository.
				GetShipmentByIDForUpdateTx(
					txCtx,
					shipmentID,
				)

			if err != nil {
				return err
			}

			switch shipment.Status {

			case model.ShipmentStatusCreated:
				// valid

			case model.ShipmentStatusCancelled:
				return ErrShipmentCancelled

			case model.ShipmentStatusDelivered:
				return ErrShipmentDelivered

			default:
				return ErrShipmentAlreadyShipped
			}

			items, err := s.repository.ListShipmentItems(
				txCtx,
				shipmentID,
			)
			if err != nil {
				return err
			}

			if len(items) == 0 {
				return ErrEmptyShipment
			}

			now := time.Now().UTC()

			for _, item := range items {

				inventory, err := s.repository.
					GetInventoryByDeviceIDForUpdateTx(
						txCtx,
						item.DeviceID,
					)

				if err != nil {
					return err
				}

				if inventory.Status != model.InventoryStatusInStock {
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

func (s *service) GetShipment(
	ctx context.Context,
	shipmentID uint,
) (*ShipmentResponse, error) {

	shipment, err := s.repository.GetShipmentByID(
		ctx,
		shipmentID,
	)
	if err != nil {
		return nil, err
	}

	items, err := s.repository.ListShipmentItems(
		ctx,
		shipmentID,
	)
	if err != nil {
		return nil, err
	}

	return &ShipmentResponse{
		ID:              shipment.ID,
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

	if req == nil {
		return ErrInvalidShipmentStatus
	}

	status := req.Status

	switch status {

	case ShipmentStatusRequest(model.ShipmentStatusCreated.String()),
		ShipmentStatusRequest(model.ShipmentStatusInTransit.String()),
		ShipmentStatusRequest(model.ShipmentStatusOutForDelivery.String()),
		ShipmentStatusRequest(model.ShipmentStatusDelivered.String()),
		ShipmentStatusRequest(model.ShipmentStatusException.String()),
		ShipmentStatusRequest(model.ShipmentStatusCancelled.String()):

	default:
		return ErrInvalidShipmentStatus
	}

	occurredAt := time.Now().UTC()

	if req.OccurredAt != nil {
		occurredAt = req.OccurredAt.UTC()
	}

	return s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {

			shipment, err := s.repository.
				GetShipmentByIDForUpdateTx(
					txCtx,
					shipmentID,
				)

			if err != nil {
				return err
			}

			externalEventID :=
				strings.TrimSpace(
					req.ExternalEventID,
				)

			// Idempotency.
			if externalEventID != "" {

				existing, err :=
					s.repository.GetTrackingEventByExternalID(
						txCtx,
						shipmentID,
						externalEventID,
					)

				if err == nil && existing != nil {
					return nil
				}

				if err != ErrTrackingEventNotFound {
					return err
				}
			}

			// Terminal shipment states cannot be downgraded.
			if shipment.Status == model.ShipmentStatusDelivered {
				return ErrShipmentDelivered
			}

			if shipment.Status == model.ShipmentStatusCancelled {
				return ErrShipmentCancelled
			}

			event := &model.ShipmentTrackingEvent{
				ShipmentID:      shipmentID,
				ExternalEventID: externalEventID,
				Status:          model.ShipmentStatus(status),
				OccurredAt:      occurredAt,
				Location: strings.TrimSpace(
					req.Location,
				),
				Description: strings.TrimSpace(
					req.Description,
				),
			}

			if err := s.repository.CreateTrackingEventTx(
				txCtx,
				event,
			); err != nil {
				return err
			}

			shipment.Status = model.ShipmentStatus(status)

			switch status {

			case ShipmentStatusRequest(model.ShipmentStatusDelivered.String()):

				deliveredAt := occurredAt
				shipment.DeliveredAt = &deliveredAt

			case ShipmentStatusRequest(model.ShipmentStatusCancelled.String()):
				// terminal state

			default:
				// keep DeliveredAt untouched
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

	if _, err := s.repository.GetShipmentByID(
		ctx,
		shipmentID,
	); err != nil {
		return nil, err
	}

	events, err := s.repository.ListTrackingEvents(
		ctx,
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
// Response Mapping
// -----------------------------------------------------------------------------

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

		result = append(
			result,
			&ShipmentItemResponse{
				ID:         item.ID,
				ShipmentID: item.ShipmentID,
				DeviceID:   item.DeviceID,
				CreatedAt:  item.CreatedAt,
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
