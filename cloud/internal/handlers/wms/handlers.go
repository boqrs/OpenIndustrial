package wms

import (
	"fmt"
	"net/http"
	"strconv"

	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/wms"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service srv.Service
}

func NewHandler(service srv.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RouterRegister(
	router ginx.ZeroGinRouter,
) {
	externalGroup := router.Group("/api/v1/external")

	externalGroup.Handle(
		http.MethodPost,
		"/warehouses",
		h.createWarehouse,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/warehouses/:warehouse_id",
		h.getWarehouse,
	)

	externalGroup.Handle(
		http.MethodPost,
		"/warehouses/:warehouse_id/locations",
		h.createLocation,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/inventory/devices/:device_id",
		h.getDeviceInventory,
	)

	externalGroup.Handle(
		http.MethodPost,
		"/inventory/inbound",
		h.stockIn,
	)

	externalGroup.Handle(
		http.MethodPost,
		"/shipments",
		h.createShipment,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/shipments/:shipment_id",
		h.getShipment,
	)

	externalGroup.Handle(
		http.MethodPost,
		"/shipments/:shipment_id/stock-out",
		h.stockOut,
	)

	externalGroup.Handle(
		http.MethodPost,
		"/shipments/:shipment_id/tracking",
		h.addTrackingEvent,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/shipments/:shipment_id/tracking",
		h.listTrackingEvents,
	)
}

func (h *Handler) createWarehouse(
	ctx *gin.Context,
) ginx.Render {

	var req srv.CreateWarehouseRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := h.service.CreateWarehouse(
		ctx.Request.Context(),
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) getWarehouse(
	ctx *gin.Context,
) ginx.Render {

	id, err := strconv.ParseUint(
		ctx.Param("warehouse_id"),
		10,
		64,
	)
	if err != nil {
		return ginx.Error(
			fmt.Errorf("invalid warehouse_id"),
		)
	}

	resp, err := h.service.GetWarehouse(
		ctx.Request.Context(),
		uint(id),
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) createLocation(
	ctx *gin.Context,
) ginx.Render {

	warehouseID, err := strconv.ParseUint(
		ctx.Param("warehouse_id"),
		10,
		64,
	)
	if err != nil {
		return ginx.Error(
			fmt.Errorf("invalid warehouse_id"),
		)
	}

	var req srv.CreateLocationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(
			fmt.Errorf("invalid param"),
		)
	}

	req.WarehouseID = uint(warehouseID)

	resp, err := h.service.CreateLocation(
		ctx.Request.Context(),
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) getDeviceInventory(
	ctx *gin.Context,
) ginx.Render {

	deviceID, err := strconv.ParseUint(
		ctx.Param("device_id"),
		10,
		64,
	)
	if err != nil {
		return ginx.Error(
			fmt.Errorf("invalid device_id"),
		)
	}

	resp, err := h.service.GetDeviceInventory(
		ctx.Request.Context(),
		uint(deviceID),
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) stockIn(
	ctx *gin.Context,
) ginx.Render {

	var req srv.StockInRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(
			fmt.Errorf("invalid param"),
		)
	}

	resp, err := h.service.StockIn(
		ctx.Request.Context(),
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) createShipment(
	ctx *gin.Context,
) ginx.Render {

	var req srv.CreateShipmentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(
			fmt.Errorf("invalid param"),
		)
	}

	resp, err := h.service.CreateShipment(
		ctx.Request.Context(),
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) getShipment(
	ctx *gin.Context,
) ginx.Render {

	shipmentID, err := strconv.ParseUint(
		ctx.Param("shipment_id"),
		10,
		64,
	)
	if err != nil {
		return ginx.Error(
			fmt.Errorf("invalid shipment_id"),
		)
	}

	resp, err := h.service.GetShipment(
		ctx.Request.Context(),
		uint(shipmentID),
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) stockOut(
	ctx *gin.Context,
) ginx.Render {

	shipmentID, err := strconv.ParseUint(
		ctx.Param("shipment_id"),
		10,
		64,
	)
	if err != nil {
		return ginx.Error(
			fmt.Errorf("invalid shipment_id"),
		)
	}

	if err := h.service.StockOut(
		ctx.Request.Context(),
		uint(shipmentID),
	); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

func (h *Handler) addTrackingEvent(
	ctx *gin.Context,
) ginx.Render {

	shipmentID, err := strconv.ParseUint(
		ctx.Param("shipment_id"),
		10,
		64,
	)
	if err != nil {
		return ginx.Error(
			fmt.Errorf("invalid shipment_id"),
		)
	}

	var req srv.TrackingEventRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(
			fmt.Errorf("invalid param"),
		)
	}

	if err := h.service.AddTrackingEvent(
		ctx.Request.Context(),
		uint(shipmentID),
		&req,
	); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

func (h *Handler) listTrackingEvents(
	ctx *gin.Context,
) ginx.Render {

	shipmentID, err := strconv.ParseUint(
		ctx.Param("shipment_id"),
		10,
		64,
	)
	if err != nil {
		return ginx.Error(
			fmt.Errorf("invalid shipment_id"),
		)
	}

	resp, err := h.service.ListTrackingEvents(
		ctx.Request.Context(),
		uint(shipmentID),
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}