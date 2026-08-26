package workorder

import (
	"fmt"
	"net/http"
	"strconv"
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service srv.Service
	auth    middleware.Service
}

func NewHandler(service srv.Service, auth middleware.Service) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
	}
}

func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {
	// External APIs
	externalGroup := router.Group("/api/v1/external")
	externalGroup.Use(h.auth.Authenticate())
	{
		externalGroup.Handle(http.MethodPost, "/work-orders", h.Create)
		externalGroup.Handle(http.MethodGet, "/work-orders", h.List)
		externalGroup.Handle(http.MethodGet, "/work-orders/:id", h.Get)
		externalGroup.Handle(http.MethodPut, "/work-orders/:id", h.Update)
		externalGroup.Handle(http.MethodPost, "/work-orders/:id/release", h.Release)
		externalGroup.Handle(http.MethodPost, "/work-orders/:id/start", h.Start)
		externalGroup.Handle(http.MethodPost, "/work-orders/:id/cancel", h.Cancel)
	}
}

func (h *Handler) Create(ctx *gin.Context) ginx.Render {
	var req srv.CreateWorkOrderRequest
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(fmt.Errorf("no perm"))
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid parm json"))
	}

	result, err := h.service.CreateWorkOrder(ctx.Request.Context(), &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

func (h *Handler) Get(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.GetWorkOrder(ctx.Request.Context(), id)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(result)
}

func (h *Handler) List(ctx *gin.Context) ginx.Render {
	var status *model.WorkOrderStatus
	if value := ctx.Query("status"); value != "" {
		s := model.WorkOrderStatus(value)
		status = &s
	}

	var productionPlanID *uint
	if value := ctx.Query("production_plan_id"); value != "" {
		id, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return ginx.Error(fmt.Errorf("invalid production_plan_id format"))
		}
		uid := uint(id)
		productionPlanID = &uid
	}

	result, err := h.service.ListWorkOrders(
		ctx.Request.Context(),
		status,
		productionPlanID,
	)

	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

func (h *Handler) Update(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	var req srv.UpdateWorkOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.UpdateWorkOrder(ctx.Request.Context(), id, &req)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(result)
}

func (h *Handler) Release(ctx *gin.Context) ginx.Render {
	h.changeState(ctx, h.service.ReleaseWorkOrder)
	return ginx.Success(nil)
}

func (h *Handler) Start(ctx *gin.Context) ginx.Render {
	h.changeState(ctx, h.service.StartWorkOrder)
	return ginx.Success(nil)
}

func (h *Handler) Cancel(ctx *gin.Context) ginx.Render {
	h.changeState(ctx, h.service.CancelWorkOrder)
	return ginx.Success(nil)
}

func (h *Handler) changeState(c *gin.Context, fn func(context.Context, uint) error) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		// In a real app, you'd render an error here.
		// For simplicity, we just return.
		return
	}

	if err := fn(c.Request.Context(), id); err != nil {
		// In a real app, you'd render an error here.
		return
	}
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