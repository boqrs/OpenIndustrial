package resource

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/boqrs/zeus/ginx"

	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"

)

// --- Permission Constants (ADDED) ---
// Using constants for permission keys improves readability and maintainability.
const (
	PermissionResourceCreate = "resource:create"
	PermissionResourceRead   = "resource:read"
	PermissionResourceUpdate = "resource:update"
	PermissionResourceDelete = "resource:delete"
)

// ResourceHandler handles HTTP requests for the resource domain.
type Handler struct {
	service srv.Service
	auth middleware.Service
}

// NewResourceHandler creates a new ResourceHandler.
// UPDATED: Removed authzRepo from parameters.
func NewHandler(service srv.Service,auth middleware.Service) *Handler {
	
	return &Handler{
		service:        service,
		auth:       auth,
	}
}

// RegisterRoutes registers the resource-related API routes.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {
	// --- Existing Product-specific Routes (PRESERVED) ---
	// Note: For simplicity, we are not applying fine-grained permissions here yet,
	// but it would follow the same pattern as the generic resources below.

	externalGroup := router.Group("/api/v1/external")
	externalGroup.Use(h.auth.Authenticate())
	externalGroup.Handle(http.MethodPost, "/products", h.handleCreateProduct)
	externalGroup.Handle(http.MethodGet, "/products", h.handleListProducts)
	externalGroup.Handle(http.MethodGet, "/products/:id", h.handleGetProduct)
	externalGroup.Handle(http.MethodPost, "/resources", h.handleCreateResource).Use(h.auth.RequirePermission(PermissionResourceCreate))
	externalGroup.Handle(http.MethodGet, "/resources", h.handleListResources).Use(h.auth.RequirePermission(PermissionResourceRead))
	externalGroup.Handle(http.MethodGet, "/resources/:id", h.handleGetResource).Use(h.auth.RequirePermission(PermissionResourceRead))
	externalGroup.Handle(http.MethodPut, "/resources/:id", h.handleUpdateResource).Use(h.auth.RequirePermission(PermissionResourceUpdate))
	externalGroup.Handle(http.MethodDelete, "/resources/:id", h.handleDeleteResource).Use(h.auth.RequirePermission(PermissionResourceDelete))
}

// --- Existing Product-specific Handlers (PRESERVED) ---
// handleCreateProduct handles the creation of a new product resource.
func (h *Handler) handleCreateProduct(ctx *gin.Context) ginx.Render {
	tenantID, err := middleware.GetTenantIDFromContext(ctx)
	if err != nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	var req srv.CreateProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}

	product, err := h.service.CreateProduct(ctx.Request.Context(), tenantID, &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(product)
}

// handleListProducts handles listing all products for the tenant.
// UPDATED: This function now supports pagination and correctly calls the new service method.
func (h *Handler) handleListProducts(ctx *gin.Context) ginx.Render {
	tenantID, err := middleware.GetTenantIDFromContext(ctx)
	if err != nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	// Read pagination parameters from query string, with defaults.
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	// In this handler, we are specifically listing resources of type "product".
	// A more generic handler might pass this type from the URL.
	resourceType := "product"

	products, err := h.service.ListResources(ctx.Request.Context(), tenantID, resourceType, limit, offset)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products: " + err.Error()})
		return ginx.Error(err)
	}

		return ginx.Success(products)
}

// handleGetProduct handles retrieving a single product by its ID.
func (h *Handler) handleGetProduct(ctx *gin.Context) ginx.Render {
	tenantID, err := middleware.GetTenantIDFromContext(ctx)
	if err != nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid credential id format"))
	}


	product, err := h.service.GetResource(ctx.Request.Context(), tenantID, uint(productID))
	if err != nil {
		
		return ginx.Error(err)
	}
	return ginx.Success(product)

}

// handleCreateResource handles creating a generic resource.
func (h *Handler) handleCreateResource(ctx *gin.Context) ginx.Render {
	tenantID, err := middleware.GetTenantIDFromContext(ctx)
	if err != nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	var req srv.CreateResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}
	req.TenantID = tenantID // Ensure the tenant ID is set from the context, not the request body.
	res, err := h.service.CreateResource(ctx.Request.Context(), &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(res)
}

// handleListResources handles listing generic resources with filtering and pagination.
func (h *Handler) handleListResources(ctx *gin.Context) ginx.Render {
	tenantID, err := middleware.GetTenantIDFromContext(ctx)
	if err != nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	resourceType := ctx.Query("type") // Allow filtering by type

	resources, err := h.service.ListResources(ctx.Request.Context(), tenantID, resourceType, limit, offset)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resources)
}

// handleGetResource handles retrieving a single generic resource by its ID.
func (h *Handler) handleGetResource(ctx *gin.Context) ginx.Render {
	tenantID, err := middleware.GetTenantIDFromContext(ctx)
	if err != nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	resourceID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid credential id format"))
	}
	

	res, err := h.service.GetResource(ctx.Request.Context(), tenantID, uint(resourceID))
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(res)
}


// handleUpdateResource handles updating a generic resource.
func (h *Handler) handleUpdateResource(ctx *gin.Context) ginx.Render {
	tenantID, err := middleware.GetTenantIDFromContext(ctx)
	if err != nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	resourceID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid credential id format"))
	}


	var req srv.UpdateResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}

	req.TenantID = tenantID
	res, err := h.service.UpdateResource(ctx.Request.Context(), uint(resourceID), &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(res)
}

// handleDeleteResource handles deleting a resource.
func (h *Handler) handleDeleteResource(ctx *gin.Context) ginx.Render {
	tenantID, err := middleware.GetTenantIDFromContext(ctx)
	if err != nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	resourceID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid credential id format"))
	}

	if err := h.service.DeleteResource(ctx.Request.Context(), tenantID, uint(resourceID)); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}


// REMOVED: handleListGroups function is removed.
// This logic belongs in an IdentityHandler that uses the Identity service.