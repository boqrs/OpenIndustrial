package api

import (
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
	// 1. Extract user and tenant info from the JWT token (set by authMiddleware)
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 2. Parse and validate the request body
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 3. Call the service layer to execute the business logic
	params := resource.CreateProductParams{
		Name:         req.Name,
		Description:  req.Description,
		Type:         req.Type,
		SerialNumber: req.SerialNumber,
		OwnerGroupID: req.OwnerGroupID,
	}

	// Correctly call the service method with tenantID as a separate argument
	newProduct, err := h.service.CreateProduct(c.Request.Context(), tenantID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product: " + err.Error()})
		return
	}

	// 4. Return the successful response
	c.JSON(http.StatusCreated, newProduct)
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