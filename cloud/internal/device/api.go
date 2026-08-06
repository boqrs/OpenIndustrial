package device

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API provides the HTTP handlers for the device domain.
type API struct {
	service *Service
}

// NewAPI creates a new device API handler.
func NewAPI(service *Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the device API routes with a Gin router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	deviceRoutes := router.Group("/devices")
	{
		deviceRoutes.POST("", a.createDevice)
		deviceRoutes.GET("", a.listDevices) // This might list all devices in an org
		deviceRoutes.GET("/:device_id", a.getDevice)
	}
}

func (a *API) createDevice(c *gin.Context) {
	orgIDStr := c.Param("org_id") // Assuming org_id is a URL parameter
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gatewayID, err := uuid.Parse(req.GatewayID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gateway id"})
		return
	}

	dev, err := a.service.CreateDevice(c.Request.Context(), orgID, gatewayID, req.Name, req.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create device"})
		return
	}

	c.JSON(http.StatusCreated, ToDeviceResponse(dev))
}

func (a *API) listDevices(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	// Check if filtering by gateway
	gatewayIDStr := c.Query("gateway_id")
	if gatewayIDStr != "" {
		gatewayID, err := uuid.Parse(gatewayIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gateway id query parameter"})
			return
		}
		devs, err := a.service.ListDevicesForGateway(c.Request.Context(), orgID, gatewayID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list devices by gateway"})
			return
		}
		c.JSON(http.StatusOK, ToDeviceListResponse(devs))
		return
	}

	devs, err := a.service.repo.ListByOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list devices"})
		return
	}

	c.JSON(http.StatusOK, ToDeviceListResponse(devs))
}

func (a *API) getDevice(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	deviceIDStr := c.Param("device_id")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	dev, err := a.service.GetDeviceByID(c.Request.Context(), orgID, deviceID)
	if err != nil {
		if err == ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get device"})
		return
	}

	c.JSON(http.StatusOK, ToDeviceResponse(dev))
}