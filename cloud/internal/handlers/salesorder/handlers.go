package salesorder

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/boqrs/zeus/ginx"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	salesorderSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/salesorder"
)

type Handler struct {
	srv  salesorderSrv.Service
	auth middleware.Service
}

func NewHandler(
	srv salesorderSrv.Service,
	auth middleware.Service,
) *Handler {
	return &Handler{
		srv:  srv,
		auth: auth,
	}
}

func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {
	externalGroup := router.Group("/api/v1/external")
	externalGroup.Use(h.auth.Authenticate())

	externalGroup.Handle(
		http.MethodPost,
		"/sales-orders",
		h.create,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/sales-orders",
		h.list,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/sales-orders/code/:code",
		h.getByOrderNo,
	)

	externalGroup.Handle(
		http.MethodGet,
		"/sales-orders/:id",
		h.getByID,
	)

	externalGroup.Handle(
		http.MethodPut,
		"/sales-orders/:id",
		h.update,
	)

	externalGroup.Handle(
		http.MethodPost,
		"/sales-orders/:id/confirm",
		h.confirm,
	)

	externalGroup.Handle(
		http.MethodPost,
		"/sales-orders/:id/cancel",
		h.cancel,
	)

}

func (h *Handler) create(ctx *gin.Context) ginx.Render {
	var req salesorderSrv.CreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)
	}

	order, err := h.srv.Create(
		ctx.Request.Context(),
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(order)

}

func (h *Handler) list(ctx *gin.Context) ginx.Render {
	var req salesorderSrv.ListRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return ginx.Error(err)
	}

	orders, err := h.srv.List(
		ctx.Request.Context(),
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(orders)

}

func (h *Handler) getByID(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(errors.New("invalid id"))
	}

	order, err := h.srv.GetByID(
		ctx.Request.Context(),
		id,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(order)

}

func (h *Handler) getByOrderNo(ctx *gin.Context) ginx.Render {
	orderNo := ctx.Param("code")

	order, err := h.srv.GetByOrderNo(
		ctx.Request.Context(),
		orderNo,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(order)

}

func (h *Handler) update(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(errors.New("invalid id"))
	}

	var req salesorderSrv.UpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)
	}

	order, err := h.srv.Update(
		ctx.Request.Context(),
		id,
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(order)

}

func (h *Handler) confirm(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(errors.New("invalid id"))
	}

	if err := h.srv.Confirm(
		ctx.Request.Context(),
		id,
	); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)

}

func (h *Handler) cancel(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(errors.New("invalid id"))
	}

	if err := h.srv.Cancel(
		ctx.Request.Context(),
		id,
	); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)

}

func parseUintParam(
	ctx *gin.Context,
	paramName string,
) (uint, error) {
	idStr := ctx.Param(paramName)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}

	return uint(id), nil

}
