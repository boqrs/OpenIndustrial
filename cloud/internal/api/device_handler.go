package api

import (
	"errors"
	"net/http"

	"github.com/OpenIndustrial/cloud/internal/device"
	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceHandler handles HTTP requests related to devices, product models, and assets.
type DeviceHandler struct {
	deviceSvc device.Service
	auth      gin.HandlerFunc // Changed: The middleware is a gin.HandlerFunc
}

// NewDeviceHandler creates a new DeviceHandler.
// It now accepts the auth middleware directly as a gin.HandlerFunc.
func NewDeviceHandler(deviceSvc device.Service, auth gin.HandlerFunc) *DeviceHandler {
	return &DeviceHandler{
		deviceSvc: deviceSvc,
		auth:      auth,
	}
}

// RegisterRoutes registers all device-related routes.
func (h *DeviceHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Product model routes
	// Create a group and then apply the middleware to it.
	productModels := router.Group("/product-models")
	productModels.Use(h.auth) // Correct way to apply middleware
	{
		productModels.POST("", h.createProductModel)
	}

	// Device instance routes
	devices := router.Group("/devices")
	devices.Use(h.auth) // Correct way to apply middleware
	{
		devices.POST("/iot-products", h.registerIoTProduct)
		devices.POST("/factory-assets", h.registerFactoryAsset)
	}
}

// createProductModel handles the API request to create a new product model.
func (h *DeviceHandler) createProductModel(c *gin.Context) {
	var req param.CreateProductModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Extract tenant ID from the context (set by auth middleware)
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found in context"})
		return
	}
	req.TenantID = tenantID.(uuid.UUID)

	// Call the service to create the product model
	productModel, err := h.deviceSvc.CreateProductModel(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product model: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, productModel)
}

// registerIoTProduct handles the API request to register a new IoT product instance.
func (h *DeviceHandler) registerIoTProduct(c *gin.Context) {
	h.registerDevice(c, device.ResourceTypeIoTProduct)
}

// registerFactoryAsset handles the API request to register a new factory asset instance.
func (h *DeviceHandler) registerFactoryAsset(c *gin.Context) {
	h.registerDevice(c, device.ResourceTypeFactoryAsset)
}

// registerDevice is a generic handler for registering any type of device instance.
func (h *DeviceHandler) registerDevice(c *gin.Context, resourceType string) {
	var req param.RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Extract tenant ID from the context
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found in context"})
		return
	}
	req.TenantID = tenantID.(uuid.UUID)

	var resp *param.ResourceResponse
	var err error

	// The service layer abstracts away the type, so we can call the appropriate method directly.
	if resourceType == device.ResourceTypeIoTProduct {
		resp, err = h.deviceSvc.RegisterIoTProduct(c.Request.Context(), &req)
	} else {
		resp, err = h.deviceSvc.RegisterFactoryAsset(c.Request.Context(), &req)
	}

	if err != nil {
		// Check for specific, known errors from the service layer
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product model not found"})
			return
		}
		// You can add more checks here for other specific errors, like invalid attributes
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register device: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}