package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IdentityHandler handles HTTP requests for the identity domain.
type IdentityHandler struct {
	service          identity.Service
	permissionRepo   PermissionRepository
	authMiddleware   gin.HandlerFunc
}

// NewIdentityHandler creates a new IdentityHandler.
func NewIdentityHandler(service identity.Service, permissionRepo PermissionRepository, authMiddleware gin.HandlerFunc) *IdentityHandler {
	return &IdentityHandler{
		service:          service,
		permissionRepo:   permissionRepo,
		authMiddleware:   authMiddleware,
	}
}

// RegisterRoutes registers the identity routes.
func (h *IdentityHandler) RegisterRoutes(router *gin.RouterGroup) {
	identityRoutes := router.Group("/identity")
	{
		identityRoutes.POST("/register", h.handleRegister)
		identityRoutes.POST("/login", h.handleLogin)
	}

	// All routes below this require authentication
	authRoutes := router.Group("/")
	authRoutes.Use(h.authMiddleware)
	{
		authRoutes.GET("/me", h.handleGetCurrentUser)

		// User management routes
		usersGroup := authRoutes.Group("/users")
		{
			usersGroup.POST("", h.handleCreateUser)
			usersGroup.GET("", h.handleListUsers)
			usersGroup.GET("/:id", h.handleGetUser)
			usersGroup.PUT("/:id", h.handleUpdateUser)
			usersGroup.DELETE("/:id", h.handleDeleteUser)
			usersGroup.POST("/:id/roles", h.handleAssignRoleToUser)
		}

		// Role management routes
		rolesGroup := authRoutes.Group("/roles")
		{
			rolesGroup.GET("", h.handleListRoles)
		}
	}
}

func (h *IdentityHandler) handleRegister(c *gin.Context) {
	var params param.RegisterTenantRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.RegisterNewTenant(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register tenant"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *IdentityHandler) handleLogin(c *gin.Context) {
	var params param.LoginRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Login(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

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

	user, err := h.service.GetCurrentUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// handleCreateUser creates a new user.
func (h *IdentityHandler) handleCreateUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var params param.CreateUserRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateUser(c.Request.Context(), tenantID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// handleListUsers lists all users in the tenant.
func (h *IdentityHandler) handleListUsers(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var params param.ListUsersRequest
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// We can bind query params here in the future, e.g., for pagination
	// if err := c.ShouldBindQuery(&params); err != nil { ... }

	users, err := h.service.ListUsers(c.Request.Context(), tenantID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// handleGetUser retrieves a single user by their ID.
func (h *IdentityHandler) handleGetUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Using GetCurrentUser because it's the same logic for now
	user, err := h.service.GetCurrentUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// handleUpdateUser updates a user.
func (h *IdentityHandler) handleUpdateUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// CORRECTED: Get userID from URL parameter
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	var params param.UpdateUserRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// CORRECTED: Pass userID from URL to the service call
	err = h.service.UpdateUser(c.Request.Context(), tenantID, userID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.Status(http.StatusOK)
}

// handleDeleteUser deletes a user.
func (h *IdentityHandler) handleDeleteUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	err = h.service.DeleteUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

// handleListRoles lists all roles.
func (h *IdentityHandler) handleListRoles(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	roles, err := h.service.ListRoles(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// handleAssignRoleToUser assigns a role to a user.
func (h *IdentityHandler) handleAssignRoleToUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// CORRECTED: Get userID from URL parameter
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	var params param.AssignRoleToUserRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// CORRECTED: Pass userID from URL to the service call
	err = h.service.AssignRoleToUser(c.Request.Context(), tenantID, userID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign role"})
		return
	}

	c.Status(http.StatusOK)
}

// --- Helper functions to get info from context ---

func getAuthPayload(c *gin.Context) (*identity.Claims, error) {
	payload, exists := c.Get("auth_payload")
	if !exists {
		return nil, errors.New("auth payload not found in context")
	}

	claims, ok := payload.(*identity.Claims)
	if !ok {
		return nil, errors.New("invalid auth payload type in context")
	}

	return claims, nil
}

// func getTenantIDFromContext(c *gin.Context) (uuid.UUID, error) {
// 	claims, err := getAuthPayload(c)
// 	if err != nil {
// 		return uuid.Nil, err
// 	}
// 	return claims.TenantID, nil
// }

// func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
// 	claims, err := getAuthPayload(c)
// 	if err != nil {
// 		return uuid.Nil, err
// 	}
// 	return claims.UserID, nil
// }