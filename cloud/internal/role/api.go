package role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API provides the HTTP handlers for the role domain.
type API struct {
	service *Service
}

// NewAPI creates a new role API handler.
func NewAPI(service *Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the role API routes with a Gin router.
// Note: These routes are typically nested under an organization route, e.g., /orgs/:org_id/roles
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	roleRoutes := router.Group("/roles")
	{
		roleRoutes.POST("", a.createRole)
		roleRoutes.GET("", a.listRoles)
		roleRoutes.GET("/:role_id", a.getRole)
	}
}

func (a *API) createRole(c *gin.Context) {
	orgIDStr := c.Param("org_id") // Assuming org_id is a URL parameter from a parent router
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := a.service.CreateRole(c.Request.Context(), orgID, req.Name, req.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, ToRoleResponse(role))
}

func (a *API) listRoles(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	roles, err := a.service.ListRolesForOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list roles"})
		return
	}

	c.JSON(http.StatusOK, ToRoleListResponse(roles))
}

func (a *API) getRole(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	roleIDStr := c.Param("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	role, err := a.service.GetRoleByID(c.Request.Context(), orgID, roleID)
	if err != nil {
		if err == ErrRoleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get role"})
		return
	}

	c.JSON(http.StatusOK, ToRoleResponse(role))
}