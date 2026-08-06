package org

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// API encapsulates all the handlers for the organization resource.
type API struct {
	service *Service
}

// NewAPI creates a new organization API handler.
func NewAPI(service *Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the organization API routes to the gin router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	orgRoutes := router.Group("/orgs")
	{
		orgRoutes.POST("", a.CreateOrganization)
		orgRoutes.GET("/:id", a.GetOrganization)
	}
}

// CreateOrganization handles the HTTP request to create a new organization.
func (a *API) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := a.service.CreateOrganization(c.Request.Context(), req.Name, req.Type, req.ParentID)
	if err != nil {
		// In a real app, you'd map domain errors to HTTP status codes.
		// For example, ErrNotFound -> 404, ErrInvalidArgument -> 400.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ToOrganizationResponse(org))
}

// GetOrganization handles the HTTP request to retrieve an organization by its ID.
func (a *API) GetOrganization(c *gin.Context) {
	id := c.Param("id")

	org, err := a.service.GetOrganization(c.Request.Context(), id)
	if err != nil {
		// Map domain errors to HTTP status codes
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToOrganizationResponse(org))
}