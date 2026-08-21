package factory

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/boqrs/zeus/ginx"
	srv"github.com/boqrs/OpenIndustrial/cloud/internal/services/factory"

)

// api handles the HTTP requests for the factory domain.
type Handler struct {
	service srv.Service
}

// NewAPI creates a new API handler for the factory service.
func NewHandler(service srv.Service) *Handler {
	return &Handler{service: service}
}

// Register registers all factory routes to the given router group.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {

	externalGroup := router.Group("/api/v1/external")
	externalGroup.Handle(http.MethodPost, "/factories", h.createFactory)
	externalGroup.Handle(http.MethodGet,"/factories/:factory_id", h.getFactory)
	externalGroup.Handle(http.MethodPut,"/factories/:factory_id", h.updateFactory)
	externalGroup.Handle(http.MethodDelete,"/factories/:factory_id", h.deleteFactory)
	externalGroup.Handle(http.MethodGet,"/factories/:factory_id/topology", h.getTopology)
	externalGroup.Handle(http.MethodPost,"/factories/topology/nodes", h.createTopologyNode)
	externalGroup.Handle(http.MethodPut,"/factories/topology/nodes/:resource_id", h.updateTopologyNode)
	externalGroup.Handle(http.MethodPost,"/factories/topology/nodes/move", h.moveTopologyNode)
	externalGroup.Handle(http.MethodDelete,"/factories/topology/nodes/:resource_id", h.deleteTopologyNode)
}

func (a *Handler) createFactory(ctx *gin.Context)ginx.Render {
	var req srv.CreateFactoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.CreateFactory(ctx.Request.Context(), &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (a *Handler) getFactory(ctx *gin.Context)ginx.Render  {
	factoryID, err := uuid.Parse(ctx.PostForm("factory_id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.GetFactory(ctx.Request.Context(), factoryID)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (a *Handler) updateFactory(ctx *gin.Context)ginx.Render  {
	factoryID, err := uuid.Parse(ctx.Param("factory_id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	var req srv.UpdateFactoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return  ginx.Error(fmt.Errorf("invalid param json"))
	}

	resp, err := a.service.UpdateFactory(ctx.Request.Context(), factoryID, &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (a *Handler) deleteFactory(ctx *gin.Context)ginx.Render  {
	factoryID, err := uuid.Parse(ctx.Param("factory_id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	if err := a.service.DeleteFactory(ctx.Request.Context(), factoryID); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)

}

func (a *Handler) createTopologyNode(ctx *gin.Context)ginx.Render  {
	var req srv.CreateTopologyNodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.CreateTopologyNode(ctx.Request.Context(), &req)
	if err != nil {
		//a.handleError(c, err)
		return ginx.Error(err)
	}

	return ginx.Success(resp)

}

func (a *Handler) updateTopologyNode(ctx *gin.Context)ginx.Render  {
	resourceID, err := uuid.Parse(ctx.Param("resource_id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	var req srv.UpdateTopologyNodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}

	resp, err := a.service.UpdateTopologyNode(ctx.Request.Context(), resourceID, &req)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(resp)
}

func (a *Handler) moveTopologyNode(ctx *gin.Context)ginx.Render  {
	var req srv.MoveTopologyNodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	if err := a.service.MoveTopologyNode(ctx.Request.Context(), &req); err != nil {
		//a.handleError(c, err)
		return ginx.Error(err)
	}

	//c.Status(http.StatusOK)
	return ginx.Success(nil)

}

func (a *Handler) deleteTopologyNode(ctx *gin.Context)ginx.Render  {
	resourceID, err := uuid.Parse(ctx.Param("resource_id"))
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource_id format"})
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	if err := a.service.DeleteTopologyNode(ctx.Request.Context(), resourceID); err != nil {
		//a.handleError(c, err)
		return ginx.Error(err)
	}

	//c.Status(http.StatusNoContent)
	return ginx.Success(nil)
}

func (a *Handler) getTopology(ctx *gin.Context)ginx.Render  {
	factoryID, err := uuid.Parse(ctx.Param("factory_id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := a.service.GetTopology(ctx.Request.Context(), factoryID)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

