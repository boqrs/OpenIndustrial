package device

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/boqrs/nexus/log"
	"github.com/boqrs/zeus/ginx"
	
)

// API handles HTTP requests for the device module.
type Handler struct {
	log    *log.Provider
	service srv.Service
}

// NewAPI creates a new API handler for the device service.
func NewHandler(service srv.Service) *Handler {
	return &Handler{service: service}
}

// Register registers all device routes to the given router group.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {

	externalGroup := router.Group("/api/v1/external")
	
	externalGroup.Handle(http.MethodGet, "/devices", h.listDevices)
	externalGroup.Handle(http.MethodGet, "/devices/:id", h.getDevice)
	externalGroup.Handle(http.MethodPatch, "/devices/:id", h.updateDevice)
	externalGroup.Handle(http.MethodDelete, "/devices/:id", h.deleteDevice)
}

//TODO: 分页逻辑后续统一定义
func (a *Handler) listDevices(ctx *gin.Context) ginx.Render {
	//page, _ := strconv.Atoi(ctx.ClientIP.DefaultQuery("page", "1"))
	req := srv.ListDevicesRequest{}
	if err := ctx.BindQuery(req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.ListDevices(ctx, &req)
	if err != nil {
		//a.handleError(c, err)
		return ginx.Error(err)
	}
	return ginx.Success(resp)
}

func (a *Handler) getDevice(ctx *gin.Context) ginx.Render {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}


	resp, err := a.service.GetDevice(ctx, uint(id))
	if err != nil {
		//a.handleError(c, err)
		return ginx.Error(err)
	}

	return ginx.Success(resp)

}

func (a *Handler) updateDevice(ctx *gin.Context) ginx.Render {
	var req srv.UpdateDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	if req.ID == nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.UpdateDevice(ctx, *req.ID, &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)

}

func (a *Handler) deleteDevice(ctx *gin.Context) ginx.Render {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	if err := a.service.DeleteDevice(ctx, uint(id)); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)

}