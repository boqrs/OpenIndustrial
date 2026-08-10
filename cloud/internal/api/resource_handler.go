package api

import (
	"database/sql"
	"net/http"

	"github.com/OpenIndustrial/cloud/internal/resource"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ResourceHandler handles HTTP requests for the resource domain.
type ResourceHandler struct {
	service         *resource.Service
	authzRepo       resource.AuthorizationRepository
	permRepo        PermissionRepository
	authMiddleware  gin.HandlerFunc
	permMiddleware  *PermissionMiddleware
}

// NewResourceHandler creates a new ResourceHandler.
func NewResourceHandler(
	service *resource.Service,
	authzRepo resource.AuthorizationRepository,
	permRepo PermissionRepository,
	authMiddleware gin.HandlerFunc,
) *ResourceHandler {
	h := &ResourceHandler{
		service:        service,
		authzRepo:      authzRepo,
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

	// Register the new route for listing groups
	groupRoutes := router.Group("/groups")
	groupRoutes.Use(h.authMiddleware)
	{
		groupRoutes.GET("", h.handleListGroups)
	}
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
func (h *ResourceHandler) handleListProducts(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	products, err := h.service.ListResources(c.Request.Context(), tenantID)
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

// handleListGroups handles listing all groups the current user is a member of.
func (h *ResourceHandler) handleListGroups(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	groups, err := h.service.ListUserGroups(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}