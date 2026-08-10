package api

import (
	"net/http"

	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IdentityHandler handles HTTP requests for the identity domain.
type IdentityHandler struct {
	service        *identity.Service
	permRepo       PermissionRepository // This now correctly refers to the single definition in permission_middleware.go
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

// RegisterRoutes registers the identity-related routes.
func (h *IdentityHandler) RegisterRoutes(router *gin.RouterGroup) {
	identityRoutes := router.Group("/identity")
	{
		identityRoutes.POST("/tenants", h.handleRegisterNewTenant)
		identityRoutes.POST("/login", h.handleLogin)
	}

	// Routes that require authentication
	authRoutes := router.Group("/identity")
	authRoutes.Use(h.authMiddleware)
	{
		authRoutes.GET("/me", h.handleGetCurrentUser)
		authRoutes.POST("/users", h.handleCreateUser)
		authRoutes.GET("/users", h.handleListUsers)
		authRoutes.PUT("/users/:userId", h.handleUpdateUser)
		authRoutes.DELETE("/users/:userId", h.handleDeleteUser)
		authRoutes.GET("/roles", h.handleListRoles)
		authRoutes.PUT("/users/:userId/roles", h.handleAssignRoleToUser)
	}
}

func (h *IdentityHandler) handleRegisterNewTenant(c *gin.Context) {
	var params identity.RegisterNewTenantParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	result, err := h.service.RegisterNewTenant(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *IdentityHandler) handleLogin(c *gin.Context) {
	var params identity.LoginParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	result, err := h.service.Login(c.Request.Context(), params)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// --- RESTORED BUSINESS LOGIC ---

func (h *IdentityHandler) handleCreateUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

func (h *IdentityHandler) handleListUsers(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

func (h *IdentityHandler) handleUpdateUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var params identity.UpdateUserParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params.UserID = userID

	if err := h.service.UpdateUser(c.Request.Context(), tenantID, params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *IdentityHandler) handleDeleteUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), tenantID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *IdentityHandler) handleListRoles(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.service.ListRoles(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *IdentityHandler) handleAssignRoleToUser(c *gin.Context) {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var params identity.AssignRoleToUserParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params.UserID = userID

	if err := h.service.AssignRoleToUser(c.Request.Context(), tenantID, params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}