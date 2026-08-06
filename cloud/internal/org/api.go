package org

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// API provides the HTTP handlers for the organization domain.
type API struct {
	service *Service
}

// NewAPI creates a new organization API handler.
func NewAPI(service *Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the organization API routes with a Gin router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	orgRoutes := router.Group("/orgs")
	{
		orgRoutes.POST("", a.createOrg)
		orgRoutes.GET("/:id", a.getOrg)
		orgRoutes.PUT("/:id", a.updateOrg)
		orgRoutes.DELETE("/:id", a.deleteOrg)
	}
}

func (a *API) createOrg(c *gin.Context) {
	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := a.service.CreateOrg(c.Request.Context(), req.Name)
	if err != nil {
		// In a real app, you'd have a proper error handling middleware
		// to map domain errors to HTTP status codes.
		if err == ErrOrgNameRequired {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create organization"})
		return
	}

	c.JSON(http.StatusCreated, ToOrgResponse(org))
}

func (a *API) getOrg(c *gin.Context) {
	id := c.Param("id")

	org, err := a.service.GetOrgByID(c.Request.Context(), id)
	if err != nil {
		if err == ErrOrgNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get organization"})
		return
	}

	c.JSON(http.StatusOK, ToOrgResponse(org))
}

func (a *API) updateOrg(c *gin.Context) {
	id := c.Param("id")
	var req UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := a.service.UpdateOrg(c.Request.Context(), id, req.Name)
	if err != nil {
		if err == ErrOrgNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == ErrOrgNameRequired {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update organization"})
		return
	}

	c.JSON(http.StatusOK, ToOrgResponse(org))
}

func (a *API) deleteOrg(c *gin.Context) {
	id := c.Param("id")

	err := a.service.DeleteOrg(c.Request.Context(), id)
	if err != nil {
		if err == ErrOrgNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete organization"})
		return
	}

	c.Status(http.StatusNoContent)
}