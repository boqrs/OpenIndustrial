package api

import (
	"net/http"

	"github.com/OpenIndustrial/cloud/internal/resource"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ResourceHandler handles HTTP requests for the Resource Kernel.
type ResourceHandler struct {
	service    *resource.Service
	authRepo   resource.AuthorizationRepository // For ABAC checks
	permRepo   PermissionRepository             // For RBAC checks (reusing from identity)
	authMiddleware gin.HandlerFunc
}

// NewResourceHandler creates a new resource handler.
func NewResourceHandler(
	service *resource.Service,
	authRepo resource.AuthorizationRepository,
	permRepo PermissionRepository,
	authMiddleware gin.HandlerFunc,
) *ResourceHandler {
	return &ResourceHandler{
		service:    service,
		authRepo:   authRepo,
		permRepo:   permRepo,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes registers the resource-related API routes.
func (h *ResourceHandler) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	products.Use(h.authMiddleware) // Secure all product routes
	{
		// We will add a permission middleware here later, e.g., h.PermissionMiddleware("products:create")
		products.POST("", h.handleCreateProduct)
	}
}

type createProductRequest struct {
	Name         string            `json:"name" binding:"required"`
	Properties   resource.Properties `json:"properties"`
	OwnerGroupID string            `json:"owner_group_id" binding:"required,uuid"`
}

func (h *ResourceHandler) handleCreateProduct(c *gin.Context) {
	// 1. Parse and validate the request body
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 2. Extract user and tenant info from the JWT token (set by authMiddleware)
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// We will need the userID for future permission checks
	// userID, _ := getUserIDFromContext(c)

	// TODO: RBAC Permission Check
	// Here we would check if the user's role has the "products:create" permission.
	// hasPermission, err := h.permRepo.CheckPermissionForUser(c.Request.Context(), userID, "products:create")
	// if err != nil || !hasPermission {
	//     c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to create products."})
	//     return
	// }

	ownerGroupID, _ := uuid.Parse(req.OwnerGroupID)

	// 3. Call the service layer to execute the business logic
	params := resource.CreateProductParams{
		TenantID:     tenantID,
		Name:         req.Name,
		Properties:   req.Properties,
		OwnerGroupID: ownerGroupID,
	}

	newProduct, err := h.service.CreateProduct(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product: " + err.Error()})
		return
	}

	// 4. Return the successful response
	c.JSON(http.StatusCreated, newProduct)
}