package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/OpenIndustrial/cloud/internal/factory"
	"github.com/OpenIndustrial/cloud/internal/param"

)

// api handles the HTTP requests for the factory domain.
type api struct {
	service factory.Service
}

// NewAPI creates a new API handler for the factory service.
func NewFactoryAPI(service factory.Service) *api {
	return &api{service: service}
}

// Register registers all factory routes to the given router group.
func (a *api) RegisterRouts(router *gin.RouterGroup) {
	factoryRoutes := router.Group("/factories")
	{
		factoryRoutes.POST("", a.createFactory)
		factoryRoutes.GET("/:factory_id", a.getFactory)
		factoryRoutes.PUT("/:factory_id", a.updateFactory)
		factoryRoutes.DELETE("/:factory_id", a.deleteFactory)
		factoryRoutes.GET("/:factory_id/topology", a.getTopology)
	}

	topologyRoutes := router.Group("/factories/topology")
	{
		topologyRoutes.POST("/nodes", a.createTopologyNode)
		topologyRoutes.PUT("/nodes/:resource_id", a.updateTopologyNode)
		topologyRoutes.POST("/nodes/move", a.moveTopologyNode)
		topologyRoutes.DELETE("/nodes/:resource_id", a.deleteTopologyNode)
	}
}

func (a *api) createFactory(c *gin.Context) {
	var req param.CreateFactoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := a.service.CreateFactory(c.Request.Context(), &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (a *api) getFactory(c *gin.Context) {
	factoryID, err := uuid.Parse(c.Param("factory_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid factory_id format"})
		return
	}

	resp, err := a.service.GetFactory(c.Request.Context(), factoryID)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (a *api) updateFactory(c *gin.Context) {
	factoryID, err := uuid.Parse(c.Param("factory_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid factory_id format"})
		return
	}

	var req param.UpdateFactoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := a.service.UpdateFactory(c.Request.Context(), factoryID, &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (a *api) deleteFactory(c *gin.Context) {
	factoryID, err := uuid.Parse(c.Param("factory_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid factory_id format"})
		return
	}

	if err := a.service.DeleteFactory(c.Request.Context(), factoryID); err != nil {
		a.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *api) createTopologyNode(c *gin.Context) {
	var req param.CreateTopologyNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := a.service.CreateTopologyNode(c.Request.Context(), &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (a *api) updateTopologyNode(c *gin.Context) {
	resourceID, err := uuid.Parse(c.Param("resource_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource_id format"})
		return
	}

	var req param.UpdateTopologyNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := a.service.UpdateTopologyNode(c.Request.Context(), resourceID, &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (a *api) moveTopologyNode(c *gin.Context) {
	var req param.MoveTopologyNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := a.service.MoveTopologyNode(c.Request.Context(), &req); err != nil {
		a.handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (a *api) deleteTopologyNode(c *gin.Context) {
	resourceID, err := uuid.Parse(c.Param("resource_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource_id format"})
		return
	}

	if err := a.service.DeleteTopologyNode(c.Request.Context(), resourceID); err != nil {
		a.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *api) getTopology(c *gin.Context) {
	factoryID, err := uuid.Parse(c.Param("factory_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid factory_id format"})
		return
	}

	resp, err := a.service.GetTopology(c.Request.Context(), factoryID)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// handleError centralizes error handling and maps domain errors to HTTP status codes.
func (a *api) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, factory.ErrFactoryNotFound), errors.Is(err, factory.ErrNodeNotFound), errors.Is(err, factory.ErrResourceNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, factory.ErrFactoryCodeExists),
		errors.Is(err, factory.ErrInvalidTopologyType),
		errors.Is(err, factory.ErrTopologyCycle),
		errors.Is(err, factory.ErrInvalidParent),
		errors.Is(err, factory.ErrCannotDeleteFactoryWithChildren),
		errors.Is(err, factory.ErrCannotDeleteNodeWithChildren),
		errors.Is(err, factory.ErrResourceTypeMismatch):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		// Log the error for internal review
		// log.Printf("internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
	}
}