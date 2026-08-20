package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/OpenIndustrial/cloud/internal/device"
	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	
)

// API handles HTTP requests for the device module.
type DeviceAPI struct {
	service device.Service
}

// NewAPI creates a new API handler for the device service.
func NewDeviceAPI(service device.Service) *DeviceAPI {
	return &DeviceAPI{service: service}
}

// Register registers all device routes to the given router group.
func (a *DeviceAPI) Register(group *gin.RouterGroup) {
	devices := group.Group("/devices")
	{
		devices.POST("", a.createDevice)
		devices.GET("", a.listDevices)
		devices.GET("/:id", a.getDevice)
		devices.PATCH("/:id", a.updateDevice)
		devices.DELETE("/:id", a.deleteDevice)
	}
}

func (a *DeviceAPI) createDevice(c *gin.Context) {
	var req param.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// Pass tenant ID from middleware into context
	// ctx := context.WithValue(c.Request.Context(), "tenant_id", getTenantIDFromGin(c))
	ctx := c.Request.Context()

	resp, err := a.service.CreateDevice(ctx, &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (a *DeviceAPI) listDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	req := param.ListDevicesRequest{
		Page:     page,
		PageSize: pageSize,
	}

	if val := c.Query("product_model_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			req.ProductModelID = &id
		}
	}
	if val := c.Query("status"); val != "" {
		status := model.DeviceStatus(val)
		req.Status = &status
	}
	if val := c.Query("parent_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			req.ParentID = &id
		}
	}

	ctx := c.Request.Context()
	resp, err := a.service.ListDevices(ctx, &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (a *DeviceAPI) getDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID format"})
		return
	}

	ctx := c.Request.Context()
	resp, err := a.service.GetDevice(ctx, id)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (a *DeviceAPI) updateDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID format"})
		return
	}

	var req param.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	resp, err := a.service.UpdateDevice(ctx, id, &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (a *DeviceAPI) deleteDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID format"})
		return
	}

	ctx := c.Request.Context()
	if err := a.service.DeleteDevice(ctx, id); err != nil {
		a.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *DeviceAPI) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, device.ErrDeviceNotFound), errors.Is(err, device.ErrProductModelNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, device.ErrSerialNumberExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, device.ErrInvalidCreateRequest), errors.Is(err, device.ErrInvalidUpdateRequest), errors.Is(err, device.ErrCannotDeleteOnlineDevice):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		// log.Printf("internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
	}
}