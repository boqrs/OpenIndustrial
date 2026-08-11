package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/OpenIndustrial/cloud/internal/kernel/resource"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ResourceHandler handles HTTP requests for the resource domain.
type ResourceHandler struct {
	service *resource.Service
	// REMOVED: authzRepo is no longer the responsibility of the resource kernel.
	// Authorization should be handled by dedicated middleware or an identity service.
	permRepo       PermissionRepository
	authMiddleware gin.HandlerFunc
	permMiddleware *PermissionMiddleware
}

// NewResourceHandler creates a new ResourceHandler.
// UPDATED: Removed authzRepo from parameters.
func NewResourceHandler(
	service *resource.Service,
	permRepo PermissionRepository,
	authMiddleware gin.HandlerFunc,
) *ResourceHandler {
	h := &ResourceHandler{
		service:        service,
		permRepo:       permRepo,
		authMiddleware: authMiddleware,
	}
	h.permMiddleware = NewPermissionMiddleware(permRepo)
	return h
}

// RegisterRoutes registers the resource-related API routes.
func (h *ResourceHandler) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	products.Use(h.authMiddleware) // Secure all product routes
	{
		products.POST("", h.handleCreateProduct)
		products.GET("", h.handleListProducts)
		products.GET("/:id", h.handleGetProduct)
	}

	// REMOVED: The /groups route is removed from this handler.
	// Group listing is an identity concern and should be handled by an IdentityHandler.
}

// CreateProductRequest defines the expected JSON body for creating a product.
type CreateProductRequest struct {
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description"`
	Type         string    `json:"type" binding:"required"`
	SerialNumber string    `json:"serial_number"`
	OwnerGroupID uuid.UUID `json:"owner_group_id" binding:"required"`
}

// handleCreateProduct handles the creation of a new product resource.
func (h *ResourceHandler) handleCreateProduct(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	params := resource.CreateProductParams{
		Name:         req.Name,
		Description:  req.Description,
		Type:         req.Type,
		SerialNumber: req.SerialNumber,
		OwnerGroupID: req.OwnerGroupID,
	}

	product, err := h.service.CreateProduct(c.Request.Context(), tenantID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// handleListProducts handles listing all products for the tenant.
// UPDATED: This function now supports pagination and correctly calls the new service method.
func (h *ResourceHandler) handleListProducts(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Read pagination parameters from query string, with defaults.
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// In this handler, we are specifically listing resources of type "product".
	// A more generic handler might pass this type from the URL.
	resourceType := "product"

	products, err := h.service.ListResources(c.Request.Context(), tenantID, resourceType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// handleGetProduct handles retrieving a single product by its ID.
func (h *ResourceHandler) handleGetProduct(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	product, err := h.service.GetResource(c.Request.Context(), tenantID, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// REMOVED: handleListGroups function is removed.
// This logic belongs in an IdentityHandler that uses the Identity service.