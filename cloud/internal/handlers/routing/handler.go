package routing

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/routing"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler manages the routing endpoints.
type Handler struct {
	service srv.Service
	auth    middleware.Service
}

// NewHandler creates a new routing handler.
func NewHandler(service srv.Service, auth middleware.Service) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
	}
}

// RouterRegister registers the routing routes.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {
	routingGroup := router.Group("/api/v1/external")
	routingGroup.Use(h.auth.Authenticate())
	{
		routingGroup.Handle(http.MethodPost, "/routings", h.CreateRouting)
		routingGroup.Handle(http.MethodGet, "/routings/:id", h.GetRouting)
		routingGroup.Handle(http.MethodPut, "/routings/:id", h.UpdateRouting)
		routingGroup.Handle(http.MethodGet, "/routings", h.ListRoutings)
		routingGroup.Handle(http.MethodPost, "/routings/:id/activate", h.ActivateRouting)
		routingGroup.Handle(http.MethodPost, "/routings/:id/deactivate", h.DeactivateRouting)
		routingGroup.Handle(http.MethodPost, "/routings/:id/operations", h.AddOperation)
		routingGroup.Handle(http.MethodGet, "/routings/:id/operations", h.ListOperations)
		routingGroup.Handle(http.MethodPut, "/routings/:id/operations/:op_id", h.UpdateOperation)
		routingGroup.Handle(http.MethodDelete, "/routings/:id/operations/:op_id", h.DeleteOperation)
	}
}

// CreateRouting creates a new routing.
func (h *Handler) CreateRouting(ctx *gin.Context) ginx.Render {
	var req srv.CreateRoutingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid request payload: %w", err))
	}

	result, err := h.service.CreateRouting(ctx.Request.Context(), &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// GetRouting retrieves a routing by its ID.
func (h *Handler) GetRouting(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.GetRouting(ctx.Request.Context(), id)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// UpdateRouting updates a routing.
func (h *Handler) UpdateRouting(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	var req srv.UpdateRoutingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid request payload: %w", err))
	}

	result, err := h.service.UpdateRouting(ctx.Request.Context(), id, &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// ActivateRouting activates a routing.
func (h *Handler) ActivateRouting(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	err = h.service.ActivateRouting(ctx.Request.Context(), id)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

// DeactivateRouting deactivates a routing.
func (h *Handler) DeactivateRouting(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	err = h.service.DeactivateRouting(ctx.Request.Context(), id)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

type listRoutingsRequest struct {
	ProductID *uint                `form:"productID"`
	Status    *model.RoutingStatus `form:"status"`
}

// ListRoutings retrieves a list of routings.
func (h *Handler) ListRoutings(ctx *gin.Context) ginx.Render {
	var req listRoutingsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid query parameters: %w", err))
	}

	result, err := h.service.ListRoutings(ctx.Request.Context(), req.ProductID, req.Status)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// AddOperation creates a new operation for a routing.
func (h *Handler) AddOperation(ctx *gin.Context) ginx.Render {
	routingID, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	var req srv.CreateOperationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid request payload: %w", err))
	}

	result, err := h.service.AddOperation(ctx.Request.Context(), routingID, &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// UpdateOperation updates an operation.
func (h *Handler) UpdateOperation(ctx *gin.Context) ginx.Render {
	routingID, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	operationID, err := uuid.Parse(ctx.Param("op_id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid operation ID format: %w", err))
	}

	var req srv.UpdateOperationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid request payload: %w", err))
	}

	result, err := h.service.UpdateOperation(ctx.Request.Context(), routingID, operationID, &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// DeleteOperation deletes an operation.
func (h *Handler) DeleteOperation(ctx *gin.Context) ginx.Render {
	routingID, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	operationID, err := uuid.Parse(ctx.Param("op_id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid operation ID format: %w", err))
	}

	err = h.service.DeleteOperation(ctx.Request.Context(), routingID, operationID)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

// ListOperations retrieves a list of operations for a routing.
func (h *Handler) ListOperations(ctx *gin.Context) ginx.Render {
	routingID, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.ListOperations(ctx.Request.Context(), routingID)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// parseUintParam is a helper to parse a uint64 from a URL parameter.
func parseUintParam(ctx *gin.Context, paramName string) (uint, error) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: '%s'", paramName, idStr)
	}
	return uint(id), nil
}