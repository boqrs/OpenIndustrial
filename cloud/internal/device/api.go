package device

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// API provides the HTTP handlers for the device domain.
type API struct {
	service Service
}

// NewAPI creates a new device API handler.
func NewAPI(service Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the device API routes.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	deviceRoutes := router.Group("/devices")
	{
		deviceRoutes.POST("", a.createDevice)
		deviceRoutes.GET("", a.listDevices)
		deviceRoutes.GET("/:device_id", a.getDevice)
	}

	gatewayRoutes := router.Group("/gateways/:gateway_id")
	{
		gatewayRoutes.GET("/devices", a.listGatewayDevices)
	}
}

func (a *API) createDevice(c *gin.Context) {
	orgID := c.Param("org_id")

	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.GatewayID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gateway_id is required"})
		return
	}

	dev, err := a.service.CreateDevice(c.Request.Context(), orgID, req.GatewayID, req.Name, req.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create device"})
		return
	}

	c.JSON(http.StatusCreated, ToDeviceResponse(dev))
}

// listDevices lists all devices for an organization.
func (a *API) listDevices(c *gin.Context) {
	orgID := c.Param("org_id")

	devices, err := a.service.ListDevicesByOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list devices"})
		return
	}
	c.JSON(http.StatusOK, ToDeviceListResponse(devices))
}

// listGatewayDevices lists all devices for a specific gateway.
func (a *API) listGatewayDevices(c *gin.Context) {
	orgID := c.Param("org_id")
	gatewayID := c.Param("gateway_id")

	devices, err := a.service.ListDevicesForGateway(c.Request.Context(), orgID, gatewayID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list devices for gateway"})
		return
	}
	c.JSON(http.StatusOK, ToDeviceListResponse(devices))
}

func (a *API) getDevice(c *gin.Context) {
	orgID := c.Param("org_id")
	deviceID := c.Param("device_id")

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