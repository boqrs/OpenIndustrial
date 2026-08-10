package api

import (
	"net/http"

	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IdentityHandler handles HTTP requests for the identity service.
type IdentityHandler struct {
	service *identity.Service
	permRepo identity.PermissionRepository // 新增
}

// NewIdentityHandler creates a new IdentityHandler.
func NewIdentityHandler(service *identity.Service, permRepo identity.PermissionRepository) *IdentityHandler {
	return &IdentityHandler{service: service,
							permRepo: permRepo,
	}
}

// RegisterRoutes registers all routes for the identity service.
func (h *IdentityHandler) RegisterRoutes(router *gin.RouterGroup) {
	identityRoutes := router.Group("/identity")
	{
		identityRoutes.POST("/tenants", h.handleRegisterNewTenant)
		identityRoutes.POST("/login", h.handleLogin)
	}

			// 受保护的路由
	authRoutes := identityRoutes.Group("/")
	authRoutes.Use(AuthMiddleware())
	{
		authRoutes.GET("/users/me", h.handleGetCurrentUser)

		authRoutes.POST("/users", h.PermissionMiddleware("identity.users:create"),h.handleCreateUser) // 新增路由
		authRoutes.GET("/users", h.PermissionMiddleware("identity.users:read"),h.handleListUsers) // 新增路由
		
		authRoutes.PUT("/users/:userId", h.PermissionMiddleware("identity.users:update"), h.handleUpdateUser) // 新增路由
		authRoutes.DELETE("/users/:userId", h.PermissionMiddleware("identity.users:delete"), h.handleDeleteUser) // 新增路由
		
		authRoutes.GET("/roles", h.PermissionMiddleware("identity.roles:read"), h.handleListRoles) // 新增路由
		authRoutes.POST("/users/:userId/roles", h.PermissionMiddleware("identity.users:assign"), h.handleAssignRoleToUser) // 新增路由
	}
}

// handleLogin is the handler for user login.
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


// handleRegisterNewTenant is the handler function for creating a new tenant.
func (h *IdentityHandler) handleRegisterNewTenant(c *gin.Context) {
	var params identity.RegisterNewTenantParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Call the service method
	result, err := h.service.RegisterNewTenant(c.Request.Context(), params)
	if err != nil {
		// In a real app, you'd have more sophisticated error handling
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// handleGetCurrentUser is the handler for getting the current user's info.
func (h *IdentityHandler) handleGetCurrentUser(c *gin.Context) {
	// 从中间件设置的 context 中获取 payload
	authPayload, err := GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.GetCurrentUser(c.Request.Context(), authPayload.TenantID, authPayload.UserID)
	if err != nil {
		// 这里可以更细致地处理错误，比如检查是否是 sql.ErrNoRows
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// handleCreateUser is the handler for creating a new user.
func (h *IdentityHandler) handleCreateUser(c *gin.Context) {
	// 从认证中间件获取租户信息
	authPayload, err := GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get auth payload"})
		return
	}

	// TODO: 在这里添加权限检查！
	// 例如: 检查 authPayload 中的用户角色是否有权限创建新用户
	// if !userHasPermission(authPayload.UserID, "create_user") {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
	// 	return
	// }

	var params identity.CreateUserParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateUser(c.Request.Context(), authPayload.TenantID, params)
	if err != nil {
		// 更好的错误处理
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// handleListUsers is the handler for listing users.
func (h *IdentityHandler) handleListUsers(c *gin.Context) {
	authPayload, err := GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get auth payload"})
		return
	}

	// TODO: Add permission check here

	var params identity.ListUsersParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ListUsers(c.Request.Context(), authPayload.TenantID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// handleListRoles is the handler for listing roles.
func (h *IdentityHandler) handleListRoles(c *gin.Context) {
	authPayload, err := GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get auth payload"})
		return
	}

	// TODO: Add permission check here (e.g., only users who can create/edit other users need to see this)

	result, err := h.service.ListRoles(c.Request.Context(), authPayload.TenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// handleAssignRoleToUser is the handler for assigning a role to a user.
func (h *IdentityHandler) handleAssignRoleToUser(c *gin.Context) {
	authPayload, err := GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get auth payload"})
		return
	}

	// 从 URL 中解析 userID
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

	// TODO: Add permission check here. Does the current user (authPayload.UserID)
	// have permission to modify roles for the target user (params.UserID)?

	err = h.service.AssignRoleToUser(c.Request.Context(), authPayload.TenantID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent) // 204 No Content is a good response for successful CUD operations without returning data.
}

// handleUpdateUser is the handler for updating a user.
func (h *IdentityHandler) handleUpdateUser(c *gin.Context) {
	authPayload, err := GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get auth payload"})
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

	// TODO: Add permission check. Can the current user update the target user?
	// (e.g., can only update self or is an admin)

	err = h.service.UpdateUser(c.Request.Context(), authPayload.TenantID, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// handleDeleteUser is the handler for deleting a user.
func (h *IdentityHandler) handleDeleteUser(c *gin.Context) {
	authPayload, err := GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get auth payload"})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	// TODO: Add permission check. Can the current user delete the target user?

	err = h.service.DeleteUser(c.Request.Context(), authPayload.TenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}