package role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API provides role-related handlers.
type API struct {
	service *Service
}

// NewAPI creates a new role API.
func NewAPI(service *Service) *API {
	return &API{service: service}
}

// RegisterRoutes registers the API routes for roles.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	// Assuming roles are nested under organizations
	orgRoutes := router.Group("/orgs/:org_id")
	{
		orgRoutes.POST("/roles", a.create)
		orgRoutes.GET("/roles", a.listForOrg)
	}
	router.GET("/roles/:id", a.getByID)
}

func (a *API) create(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("org_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := a.service.CreateRole(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ToRoleResponse(role))
}

func (a *API) listForOrg(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("org_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	roles, err := a.service.ListRolesForOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ToRoleListResponse(roles))
}

func (a *API) getByID(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := a.service.GetRoleByID(c.Request.Context(), roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, ToRoleResponse(role))
}