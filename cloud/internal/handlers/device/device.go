package device

import (
	"fmt"
	"net/http"
	"strconv"

	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/boqrs/nexus/log"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
)

// API handles HTTP requests for the device module.
type Handler struct {
	log     *log.Provider
	service srv.Service
}

// NewAPI creates a new API handler for the device service.
func NewHandler(service srv.Service) *Handler {
	return &Handler{service: service}
}

// Register registers all device routes to the given router group.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {

	externalGroup := router.Group("/api/v1/external")

	externalGroup.Handle(
		http.MethodGet,
		"/devices",
		h.listDevices,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/devices/:id",
		h.getDevice,
	)

	externalGroup.Handle(
		http.MethodPatch,
		"/devices/:id",
		h.updateDevice,
	)

	// Device activation is an end-user operation.
	//
	// The authenticated user is obtained from the JWT context
	// inside DeviceService.ActivateDevice. The request body must
	// therefore never contain customer/user ID.
	externalGroup.Handle(
		http.MethodPost,
		"/devices/activate",
		h.activateDevice,
	)

	// externalGroup.Handle(http.MethodDelete, "/devices/:id", h.deleteDevice)
}

// TODO: 分页逻辑后续统一定义
func (a *Handler) listDevices(ctx *gin.Context) ginx.Render {
	req := srv.ListDevicesRequest{}

	if err := ctx.BindQuery(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.ListDevices(
		ctx,
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (a *Handler) getDevice(ctx *gin.Context) ginx.Render {
	id, err := strconv.ParseUint(
		ctx.Param("id"),
		10,
		32,
	)
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.GetDevice(
		ctx,
		uint(id),
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (a *Handler) updateDevice(ctx *gin.Context) ginx.Render {
	var req srv.UpdateDeviceRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	if req.ID == nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.UpdateDevice(
		ctx,
		*req.ID,
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

// activateDevice activates a manufactured device for the authenticated
// end user.
//
// The request contains only the device information encoded in the QR
// code, for example:
//
//	{
//	    "serial_number": "SN202609220001",
//	    "product_model": "OI-PRINTER-X1"
//	}
//
// The customer/user UUID is intentionally NOT accepted here.
// DeviceService obtains it from the authenticated JWT context.
func (a *Handler) activateDevice(ctx *gin.Context) ginx.Render {
	var req srv.ActivateDeviceRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.ActivateDevice(
		ctx,
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

// func (a *Handler) deleteDevice(ctx *gin.Context) ginx.Render {
// 	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
// 	if err != nil {
// 		return ginx.Error(fmt.Errorf("invalid param"))
// 	}
//
// 	if err := a.service.DeleteDevice(ctx, uint(id)); err != nil {
// 		return ginx.Error(err)
// 	}
//
// 	return ginx.Success(nil)
// }
