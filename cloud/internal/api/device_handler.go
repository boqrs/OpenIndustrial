package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/OpenIndustrial/cloud/internal/device"
	"github.com/OpenIndustrial/cloud/internal/param"

)

// handler handles HTTP requests for the device domain.
type handler struct {
	svc device.Service
}

func NewDeviceAPI(service device.Service) *handler {
	return &handler{svc: service}
}

// RegisterRoutes registers all device API routes to the given router group.
func (h *handler)RegisterRoutes(rg *gin.RouterGroup) {
	dt := rg.Group("/device-types")
	{
		dt.POST("", h.createDeviceType)
		dt.GET("", h.listDeviceTypes)
		dt.GET("/:id", h.getDeviceType)
		dt.PUT("/:id", h.updateDeviceType)
	}

	// Device routes
	d := rg.Group("/devices")
	{
		d.POST("", h.createDevice)
		d.GET("", h.listDevices)
		d.GET("/:id", h.getDevice)
		d.PUT("/:id", h.updateDevice)
		d.DELETE("/:id", h.deleteDevice)

		// Topology routes
		d.POST("/:id/attach", h.attachDevice)
		d.POST("/:id/detach", h.detachDevice)
		d.POST("/:id/connect", h.connectDevice)
		d.POST("/:id/disconnect/:connectionId", h.disconnectDevice)
		d.GET("/:id/topology", h.getTopology)
	}
}

// --- DeviceType Handlers ---

func (h *handler) createDeviceType(c *gin.Context) {
	var req param.CreateDeviceTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.CreateDeviceType(c.Request.Context(), &req)
	if err != nil {
		// Consider mapping specific errors to status codes
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *handler) listDeviceTypes(c *gin.Context) {
	resp, err := h.svc.ListDeviceTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) getDeviceType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device type ID"})
		return
	}

	resp, err := h.svc.GetDeviceType(c.Request.Context(), id)
	if err != nil {
		// if err == ErrDeviceTypeNotFound {
		// 	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 	return
		// }
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) updateDeviceType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device type ID"})
		return
	}

	var req param.UpdateDeviceTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.UpdateDeviceType(c.Request.Context(), id, &req)
	if err != nil {
		// if err == ErrDeviceTypeNotFound {
		// 	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 	return
		// }
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// --- Device Handlers ---

func (h *handler) createDevice(c *gin.Context) {
	var req param.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.CreateDevice(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *handler) listDevices(c *gin.Context) {
	var deviceTypeID *uuid.UUID
	if dtIDStr := c.Query("device_type_id"); dtIDStr != "" {
		id, err := uuid.Parse(dtIDStr)
		if err == nil {
			deviceTypeID = &id
		}
	}

	resp, err := h.svc.ListDevices(c.Request.Context(), deviceTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) getDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	resp, err := h.svc.GetDevice(c.Request.Context(), id)
	if err != nil {
		// if err == ErrDeviceNotFound {
		// 	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 	return
		// }
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) updateDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	var req param.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.UpdateDevice(c.Request.Context(), id, &req)
	if err != nil {
		// if err == ErrDeviceNotFound {
		// 	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 	return
		// }
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) deleteDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	err = h.svc.DeleteDevice(c.Request.Context(), id)
	if err != nil {
		// if err == ErrDeviceNotFound {
		// 	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 	return
		// }
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// --- Topology Handlers ---

func (h *handler) attachDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	var req param.AttachDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.AttachDevice(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *handler) detachDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	err = h.svc.DetachDevice(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *handler) connectDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	var req param.ConnectDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.ConnectDevice(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *handler) disconnectDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	connID, err := strconv.ParseUint(c.Param("connectionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid connection ID"})
		return
	}

	err = h.svc.DisconnectDevice(c.Request.Context(), id, uint(connID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *handler) getTopology(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	resp, err := h.svc.GetTopology(c.Request.Context(), id)
	if err != nil {
		// if err == ErrDeviceNotFound {
		// 	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 	return
		// }
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}