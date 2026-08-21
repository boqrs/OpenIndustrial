package device

import (
	"errors"
	"net/http"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
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
	
	externalGroup.Handle(http.MethodPost, "", h.createDevice)
	externalGroup.Handle(http.MethodGet, "", h.listDevices)
	externalGroup.Handle(http.MethodGet, "/:id", h.getDevice)
	externalGroup.Handle(http.MethodPatch, "/:id", h.updateDevice)
	externalGroup.Handle(http.MethodDelete, "/:id", h.deleteDevice)
}

func (h *Handler) createDevice(ctx *gin.Context) ginx.Render {
	var req srv.CreateDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}
	// Pass tenant ID from middleware into context
	// ctx := context.WithValue(c.Request.Context(), "tenant_id", getTenantIDFromGin(c))
	resp, err := h.service.CreateDevice(ctx, &req)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(resp)
}

//TODO: 分页逻辑后续统一定义
func (a *Handler) listDevices(ctx *gin.Context) ginx.Render {
	//page, _ := strconv.Atoi(ctx.ClientIP.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	req := srv.ListDevicesRequest{
		Page:     1,
		PageSize: pageSize,
	}

	if val := ctx.Query("product_model_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			req.ProductModelID = &id
		}
	}
	if val := ctx.Query("status"); val != "" {
		status := model.DeviceStatus(val)
		req.Status = &status
	}
	if val := ctx.Query("parent_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			req.ParentID = &id
		}
	}

	resp, err := a.service.ListDevices(ctx, &req)
	if err != nil {
		//a.handleError(c, err)
		return ginx.Error(err)
	}
	return ginx.Success(resp)
}

func (a *Handler) getDevice(ctx *gin.Context) ginx.Render {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID format"})
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.GetDevice(ctx, id)
	if err != nil {
		//a.handleError(c, err)
		return ginx.Error(err)
	}

	return ginx.Success(resp)

}

func (a *Handler) updateDevice(ctx *gin.Context) ginx.Render {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	var req srv.UpdateDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.UpdateDevice(ctx, id, &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)

}

func (a *Handler) deleteDevice(ctx *gin.Context) ginx.Render {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	if err := a.service.DeleteDevice(ctx, id); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)

}

func (a *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, srv.ErrDeviceNotFound), errors.Is(err, srv.ErrProductModelNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, srv.ErrSerialNumberExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, srv.ErrInvalidCreateRequest), errors.Is(err, srv.ErrInvalidUpdateRequest), errors.Is(err, srv.ErrCannotDeleteOnlineDevice):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		// log.Printf("internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
	}
}