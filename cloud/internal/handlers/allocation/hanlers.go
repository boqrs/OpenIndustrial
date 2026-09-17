package allocation

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/boqrs/zeus/ginx"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	allocationService "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/allocation"
)

type Handler struct {
	srv  allocationService.Service
	auth middleware.Service
}

func NewHandler(
	srv allocationService.Service,
	auth middleware.Service,
) *Handler {
	return &Handler{
		srv:  srv,
		auth: auth,
	}
}

func (h *Handler) RouterRegister(
	router ginx.ZeroGinRouter,
) {
	externalGroup := router.Group("/api/v1/external")
	externalGroup.Use(h.auth.Authenticate())

	externalGroup.Handle(
		http.MethodPost,
		"/production-plan-allocations",
		h.create,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/production-plan-allocations/:id",
		h.getByID,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/sales-order-items/:id/production-plan-allocations",
		h.listBySalesOrderItemID,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/production-plans/:id/allocations",
		h.listByProductionPlanID,
	)
}

func (h *Handler) create(
	ctx *gin.Context,
) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}

	var req allocationService.CreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)
	}

	result, err := h.srv.Create(
		ctx.Request.Context(),
		tenantID,
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

func (h *Handler) getByID(
	ctx *gin.Context,
) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}

	id, err := parseUintParam(ctx, "id")
	if err != nil || id == 0 {
		return ginx.Error(errors.New("invalid id"))
	}

	result, err := h.srv.GetByID(
		ctx.Request.Context(),
		tenantID,
		id,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

func (h *Handler) listBySalesOrderItemID(
	ctx *gin.Context,
) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}

	id, err := parseUintParam(ctx, "id")
	if err != nil || id == 0 {
		return ginx.Error(errors.New("invalid id"))
	}

	results, err := h.srv.ListBySalesOrderItemID(
		ctx.Request.Context(),
		tenantID,
		id,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(results)
}

func (h *Handler) listByProductionPlanID(
	ctx *gin.Context,
) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}

	id, err := parseUintParam(ctx, "id")
	if err != nil || id == 0 {
		return ginx.Error(errors.New("invalid id"))
	}

	results, err := h.srv.ListByProductionPlanID(
		ctx.Request.Context(),
		tenantID,
		id,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(results)
}

func parseUintParam(
	ctx *gin.Context,
	paramName string,
) (uint, error) {
	idStr := ctx.Param(paramName)

	id, err := strconv.ParseUint(
		idStr,
		10,
		32,
	)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
