package device

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API provides device-related endpoints.
type API struct {
	service *Service
}

// NewAPI creates a new device API handler.
func NewAPI(service *Service) *API {
	return &API{service: service}
}

// RegisterRoutes registers device routes.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	deviceRoutes := router.Group("/devices")
	{
		deviceRoutes.POST("", a.registerDevice)
		// Other routes can be added here later.
	}
}

// registerDevice handles the creation of a new device.
func (a *API) registerDevice(c *gin.Context) {
	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orgID := uuid.New() // Placeholder

	device, err := a.service.RegisterDevice(c.Request.Context(), orgID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register device"})
		return
	}

	c.JSON(http.StatusCreated, ToDeviceResponse(device))
}