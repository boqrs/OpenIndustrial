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
	// --- Existing Product-specific Routes (PRESERVED) ---
	products := router.Group("/products")
	products.Use(h.authMiddleware) // Secure all product routes
	{
		products.POST("", h.handleCreateProduct)
		products.GET("", h.handleListProducts)
		products.GET("/:id", h.handleGetProduct)
	}

	// --- New Generic Resource Routes (ADDED) ---
	resources := router.Group("/resources")
	resources.Use(h.authMiddleware)
	{
		resources.POST("", h.handleCreateResource)
		resources.GET("", h.handleListResources)
		resources.GET("/:id", h.handleGetResource)
		resources.PUT("/:id", h.handleUpdateResource)
		resources.DELETE("/:id", h.handleDeleteResource)
	}

	// REMOVED: The /groups route is removed from this handler.
	// Group listing is an identity concern and should be handled by an IdentityHandler.
}

// --- Existing Product-specific Handlers (PRESERVED) ---

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

// --- New Generic Resource Handlers (ADDED) ---

// CreateResourceRequest defines the body for creating a generic resource.
type CreateResourceRequest struct {
	Type         string     `json:"type" binding:"required"`
	Name         string     `json:"name" binding:"required"`
	Code         *string    `json:"code"`
	Status       string     `json:"status"`
	Metadata     []byte     `json:"metadata"`
	ParentID     *uuid.UUID `json:"parent_id"`
	OwnerGroupID *uuid.UUID `json:"owner_group_id"`
}

// handleCreateResource handles creating a generic resource.
func (h *ResourceHandler) handleCreateResource(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	params := resource.CreateResourceParams{
		TenantID:     tenantID,
		Type:         req.Type,
		Name:         req.Name,
		Code:         req.Code,
		Status:       req.Status,
		Metadata:     req.Metadata,
		ParentID:     req.ParentID,
		OwnerGroupID: req.OwnerGroupID,
	}

	res, err := h.service.CreateResource(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resource: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// handleListResources handles listing generic resources with filtering and pagination.
func (h *ResourceHandler) handleListResources(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	resourceType := c.Query("type") // Allow filtering by type

	resources, err := h.service.ListResources(c.Request.Context(), tenantID, resourceType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list resources: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resources)
}

// handleGetResource handles retrieving a single generic resource by its ID.
func (h *ResourceHandler) handleGetResource(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resourceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID format"})
		return
	}

	res, err := h.service.GetResource(c.Request.Context(), tenantID, resourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get resource: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateResourceRequest defines the body for updating a resource.
type UpdateResourceRequest struct {
	Name     string `json:"name" binding:"required"`
	Code     *string `json:"code"`
	Status   string `json:"status"`
	Metadata []byte `json:"metadata"`
	Version  int    `json:"version" binding:"required"` // For optimistic locking
}

// handleUpdateResource handles updating a generic resource.
func (h *ResourceHandler) handleUpdateResource(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resourceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID format"})
		return
	}

	var req UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	params := resource.UpdateResourceParams{
		Name:     req.Name,
		Code:     req.Code,
		Status:   req.Status,
		Metadata: req.Metadata,
	}

	res, err := h.service.UpdateResource(c.Request.Context(), tenantID, resourceID, req.Version, params)
	if err != nil {
		// Specific error handling for optimistic locking
		if err.Error() == "update conflict: resource has been modified by another process" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resource: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// handleDeleteResource handles deleting a resource.
func (h *ResourceHandler) handleDeleteResource(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resourceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID format"})
		return
	}

	if err := h.service.DeleteResource(c.Request.Context(), tenantID, resourceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}


// REMOVED: handleListGroups function is removed.
// This logic belongs in an IdentityHandler that uses the Identity service.