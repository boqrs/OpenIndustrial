package planning

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/planning"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
)

// Handler manages production plan endpoints.
type Handler struct {
	service srv.Service
	auth    middleware.Service
}

// NewHandler creates a new production plan handler.
func NewHandler(
	service srv.Service,
	auth middleware.Service,
) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
	}
}

// RouterRegister registers production plan routes.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {
	group := router.Group("/api/v1/external")
	group.Use(h.auth.Authenticate())

	group.Handle(
		http.MethodPost,
		"/production-plans",
		h.CreateProductionPlan,
	)

	group.Handle(
		http.MethodGet,
		"/production-plans",
		h.ListProductionPlans,
	)

	group.Handle(
		http.MethodGet,
		"/production-plans/:id",
		h.GetProductionPlan,
	)

	group.Handle(
		http.MethodPut,
		"/production-plans/:id",
		h.UpdateProductionPlan,
	)

	group.Handle(
		http.MethodPost,
		"/production-plans/:id/release",
		h.ReleaseProductionPlan,
	)

	group.Handle(
		http.MethodPost,
		"/production-plans/:id/cancel",
		h.CancelProductionPlan,
	)
}

// CreateProductionPlan creates a production plan.
func (h *Handler) CreateProductionPlan(ctx *gin.Context) ginx.Render {
	var req srv.CreateProductionPlanRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(
			fmt.Errorf("invalid request payload: %w", err),
		)
	}

	result, err := h.service.CreateProductionPlan(
		ctx.Request.Context(),
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// GetProductionPlan retrieves a production plan by ID.
func (h *Handler) GetProductionPlan(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.GetProductionPlanByID(
		ctx.Request.Context(),
		id,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

type listProductionPlansRequest struct {
	Status *model.ProductionPlanStatus `form:"status"`
}

// ListProductionPlans lists production plans.
func (h *Handler) ListProductionPlans(ctx *gin.Context) ginx.Render {
	var req listProductionPlansRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return ginx.Error(
			fmt.Errorf("invalid query parameters: %w", err),
		)
	}

	result, err := h.service.ListProductionPlans(
		ctx.Request.Context(),
		req.Status,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// UpdateProductionPlan updates a draft production plan.
func (h *Handler) UpdateProductionPlan(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	var req srv.UpdateProductionPlanRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(
			fmt.Errorf("invalid request payload: %w", err),
		)
	}

	result, err := h.service.UpdateProductionPlan(
		ctx.Request.Context(),
		id,
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// ReleaseProductionPlan releases a draft production plan.
func (h *Handler) ReleaseProductionPlan(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	if err := h.service.ReleaseProductionPlan(
		ctx.Request.Context(),
		id,
	); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

// CancelProductionPlan cancels a production plan.
func (h *Handler) CancelProductionPlan(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	if err := h.service.CancelProductionPlan(
		ctx.Request.Context(),
		id,
	); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

func parseUintParam(ctx *gin.Context, paramName string) (uint, error) {
	idStr := ctx.Param(paramName)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf(
			"invalid %s format: '%s'",
			paramName,
			idStr,
		)
	}

	return uint(id), nil
}
