package permission

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// API provides the HTTP handlers for the permission domain.
type API struct {
	service *Service
}

// NewAPI creates a new permission API handler.
func NewAPI(service *Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the permission API routes with a Gin router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	permissionRoutes := router.Group("/permissions")
	{
		permissionRoutes.GET("", a.listPermissions)
	}
}

func (a *API) listPermissions(c *gin.Context) {
	permissions, err := a.service.ListAllPermissions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list permissions"})
		return
	}

	c.JSON(http.StatusOK, ToPermissionListResponse(permissions))
}