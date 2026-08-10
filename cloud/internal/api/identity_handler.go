package api

import (
	"context"
	"net/http"

	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PermissionRepository is an interface for checking permissions.
type PermissionRepository interface {
	CheckPermissionForUser(ctx context.Context, userID uuid.UUID, permissionName string) (bool, error)
}

// IdentityHandler handles HTTP requests for the identity service.
type IdentityHandler struct {
	service        *identity.Service
	permRepo       PermissionRepository
	authMiddleware gin.HandlerFunc
}

// NewIdentityHandler creates a new IdentityHandler.
func NewIdentityHandler(service *identity.Service, permRepo PermissionRepository, authMiddleware gin.HandlerFunc) *IdentityHandler {
	return &IdentityHandler{
		service:        service,
		permRepo:       permRepo,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes registers all routes for the identity service.
func (h *IdentityHandler) RegisterRoutes(router *gin.RouterGroup) {
	identityRoutes := router.Group("/identity")
	{
		identityRoutes.POST("/tenants", h.handleRegisterNewTenant)
		identityRoutes.POST("/login", h.handleLogin)

		authRoutes := identityRoutes.Group("/")
		authRoutes.Use(h.authMiddleware)
		{
			authRoutes.GET("/users/me", h.handleGetCurrentUser)
			authRoutes.POST("/users", h.PermissionMiddleware("identity.users:create"), h.handleCreateUser)
			authRoutes.GET("/users", h.PermissionMiddleware("identity.users:read"), h.handleListUsers)
			authRoutes.PUT("/users/:userId", h.PermissionMiddleware("identity.users:update"), h.handleUpdateUser)
			authRoutes.DELETE("/users/:userId", h.PermissionMiddleware("identity.users:delete"), h.handleDeleteUser)
			authRoutes.GET("/roles", h.PermissionMiddleware("identity.roles:read"), h.handleListRoles)
			authRoutes.POST("/users/:userId/roles", h.PermissionMiddleware("identity.users:assign"), h.handleAssignRoleToUser)
		}
	}
}

// handleRegisterNewTenant is the handler function for creating a new tenant.
func (h *IdentityHandler) handleRegisterNewTenant(c *gin.Context) {
	// CORRECTED: Using RegisterNewTenantParams as per your original code.
	var params identity.RegisterNewTenantParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// CORRECTED: Matching the 2 return values (result, error).
	result, err := h.service.RegisterNewTenant(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// handleLogin is the handler for user login.
func (h *IdentityHandler) handleLogin(c *gin.Context) {
	// CORRECTED: Using LoginParams struct as required by the service.
	var params identity.LoginParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// CORRECTED: Passing the params struct, not individual fields.
	result, err := h.service.Login(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// handleGetCurrentUser is the handler for getting the current user's info.
func (h *IdentityHandler) handleGetCurrentUser(c *gin.Context) {
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

	// CORRECTED: Using GetCurrentUser instead of the non-existent GetUserByID.
	user, err := h.service.GetCurrentUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// handleCreateUser is the handler for creating a new user.
func (h *IdentityHandler) handleCreateUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var params identity.CreateUserParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateUser(c.Request.Context(), tenantID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// handleListUsers is the handler for listing users.
func (h *IdentityHandler) handleListUsers(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var params identity.ListUsersParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ListUsers(c.Request.Context(), tenantID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// handleListRoles is the handler for listing roles.
func (h *IdentityHandler) handleListRoles(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ListRoles(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// handleAssignRoleToUser is the handler for assigning a role to a user.
func (h *IdentityHandler) handleAssignRoleToUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	var reqBody struct {
		RoleID uuid.UUID `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := identity.AssignRoleToUserParams{
		UserID: userID,
		RoleID: reqBody.RoleID,
	}

	err = h.service.AssignRoleToUser(c.Request.Context(), tenantID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// handleUpdateUser is the handler for updating a user.
func (h *IdentityHandler) handleUpdateUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	var reqBody identity.UpdateUserParams
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reqBody.UserID = userID

	err = h.service.UpdateUser(c.Request.Context(), tenantID, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// handleDeleteUser is the handler for deleting a user.
func (h *IdentityHandler) handleDeleteUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	err = h.service.DeleteUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// PermissionMiddleware creates a middleware for checking a specific permission.
func (h *IdentityHandler) PermissionMiddleware(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserIDFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		hasPerm, err := h.permRepo.CheckPermissionForUser(c.Request.Context(), userID, permissionName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			return
		}
		if !hasPerm {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have the required permission"})
			return
		}
		c.Next()
	}
}