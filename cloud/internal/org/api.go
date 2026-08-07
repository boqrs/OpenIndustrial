package org

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API provides organization-related handlers.
type API struct {
	service *Service
}

// NewAPI creates a new organization API.
func NewAPI(service *Service) *API {
	return &API{service: service}
}

// RegisterRoutes registers the API routes for organizations.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	orgRoutes := router.Group("/orgs")
	{
		orgRoutes.POST("", a.create)
		orgRoutes.GET("/:id", a.getByID)
	}
}

func (a *API) create(c *gin.Context) {
	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := a.service.CreateOrganization(c.Request.Context(), req.Name, req.Description, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ToOrgResponse(org))
}

func (a *API) getByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	org, err := a.service.GetOrganization(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, ToOrgResponse(org))
}