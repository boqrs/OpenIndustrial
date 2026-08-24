package identity

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"


	"github.com/boqrs/zeus/ginx"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/identity"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"

)

// IdentityHandler handles HTTP requests for the identity domain.
type Handler struct {
	service          srv.Service
	//permissionRepo   PermissionRepository
	auth middleware.Service
}

// NewIdentityHandler creates a new IdentityHandler.
func NewIdentityHandler(service srv.Service, auth middleware.Service) *Handler {
	return &Handler{
		service:          service,
		auth: 			  auth,
	}
}

// RegisterRoutes registers the identity routes.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {
	
	externalGroup := router.Group("/api/v1/external")
	externalGroup.Use(h.auth.Authenticate())
	externalGroup.Handle(http.MethodPost, "/identity/register", h.handleRegister)
	externalGroup.Handle(http.MethodPost, "/identity/login", h.handleLogin)
	externalGroup.Handle(http.MethodGet, "/identity/me", h.handleGetCurrentUser)
	externalGroup.Handle(http.MethodPost, "/users", h.handleCreateUser)
	externalGroup.Handle(http.MethodGet, "/users", h.handleListUsers)
	externalGroup.Handle(http.MethodGet, "/users/:id", h.handleGetUser)
	externalGroup.Handle(http.MethodPut, "/users/:id", h.handleUpdateUser)
	externalGroup.Handle(http.MethodDelete, "/users/:id", h.handleDeleteUser)
	externalGroup.Handle(http.MethodPost, "/:id/roles", h.handleAssignRoleToUser)
}

func (h *Handler) handleRegister(ctx *gin.Context) ginx.Render {
	var params srv.RegisterTenantRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	result, err := h.service.RegisterNewTenant(ctx.Request.Context(), &params)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)

}

func (h *Handler) handleLogin(ctx *gin.Context) ginx.Render {
	var params srv.LoginRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	result, err := h.service.Login(ctx.Request.Context(), &params)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

func (h *Handler) handleGetCurrentUser(ctx *gin.Context) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	userID := pkg.GetUserIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	user, err := h.service.GetCurrentUser(ctx.Request.Context(), tenantID, userID)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(user)
}

// handleCreateUser creates a new user.
func (h *Handler) handleCreateUser(ctx *gin.Context) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	var params srv.CreateUserRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	result, err := h.service.CreateUser(ctx.Request.Context(), tenantID, &params)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)

}

// handleListUsers lists all users in the tenant.
func (h *Handler) handleListUsers(ctx *gin.Context) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	var params srv.ListUsersRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	users, err := h.service.ListUsers(ctx.Request.Context(), tenantID, &params)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(users)
}

// handleGetUser retrieves a single user by their ID.
func (h *Handler) handleGetUser(ctx *gin.Context) ginx.Render{
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return ginx.Error(err)
	}

	// Using GetCurrentUser because it's the same logic for now
	user, err := h.service.GetCurrentUser(ctx.Request.Context(), tenantID, userID)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(user)

}

// handleUpdateUser updates a user.
func (h *Handler) handleUpdateUser(ctx *gin.Context) ginx.Render  {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ginx.Error(err)
	}

	var params srv.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	// CORRECTED: Pass userID from URL to the service call
	err = h.service.UpdateUser(ctx.Request.Context(), tenantID, userID, &params)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

// handleDeleteUser deletes a user.
func (h *Handler) handleDeleteUser(ctx *gin.Context) ginx.Render  {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ginx.Error(err)
	}

	if err = h.service.DeleteUser(ctx.Request.Context(), tenantID, userID); err != nil {
		return ginx.Error(err)
	}
	
	return ginx.Success(nil)
}

// handleListRoles lists all roles.
func (h *Handler) handleListRoles(ctx *gin.Context) ginx.Render  {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	roles, err := h.service.ListRoles(ctx.Request.Context(), tenantID)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(roles)

}

// handleAssignRoleToUser assigns a role to a user.
func (h *Handler) handleAssignRoleToUser(ctx *gin.Context) ginx.Render  {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	// CORRECTED: Get userID from URL parameter
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ginx.Error(err)
	}

	var params srv.AssignRoleToUserRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	// CORRECTED: Pass userID from URL to the service call
	if err = h.service.AssignRoleToUser(ctx.Copy().Request.Context(), tenantID, userID, &params);err != nil {
		return ginx.Error(err)
	}

	//c.Status(http.StatusOK)
	return ginx.Success(nil)
}

// --- Helper functions to get info from context ---

func getAuthPayload(c *gin.Context) (*srv.Claims, error) {
	payload, exists := c.Get("auth_payload")
	if !exists {
		return nil, errors.New("auth payload not found in context")
	}

	claims, ok := payload.(*srv.Claims)
	if !ok {
		return nil, errors.New("invalid auth payload type in context")
	}

	return claims, nil
}
