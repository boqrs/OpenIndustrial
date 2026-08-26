package material

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/boqrs/zeus/ginx"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	mSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/material"
)

type Handler struct {
	srv  mSrv.Service
	auth middleware.Service
}

func NewHandler(srv mSrv.Service, auth middleware.Service) *Handler {
	return &Handler{
		srv:  srv,
		auth: auth,
	}
}

func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {

	externalGroup := router.Group("/api/v1/external")
	externalGroup.Use(h.auth.Authenticate())

	externalGroup.Handle(http.MethodPost, "/materials", h.create)
	externalGroup.Handle(http.MethodGet, "/materials", h.list)
	externalGroup.Handle(http.MethodGet, "/materials/:id", h.getByID)
	externalGroup.Handle(http.MethodPut, "/materials/:id", h.update)
	externalGroup.Handle(http.MethodDelete, "/materials/:id", h.delete)

}

func (h *Handler) create(ctx *gin.Context) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID ==  uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}

	var req model.Material
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)
	}

	if err := h.srv.Create(ctx.Request.Context(), tenantID, &req); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

func (h *Handler) list(ctx *gin.Context) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}

	offsetStr := ctx.DefaultQuery("offset", "0")
	limitStr := ctx.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // 如果转换失败或值无效，使用默认值
	}

	materials, _, err := h.srv.List(ctx.Request.Context(), tenantID, offset, limit)
	if err != nil {
		return ginx.Error(err)
	}

	// 注意：这里我们使用 ginx.Page 来返回带分页信息的结果
	return ginx.Success(materials)
}

func (h *Handler) getByID(ctx *gin.Context) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID ==  uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}


	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(errors.New("invalid parma"))
	}

	material, err := h.srv.GetByID(ctx.Request.Context(), tenantID, id)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(material)
}

func (h *Handler) update(ctx *gin.Context) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID ==  uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}

	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(errors.ErrUnsupported)
	}

	var req model.Material
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)
	}
	req.ID = id

	if err := h.srv.Update(ctx.Request.Context(), tenantID, &req); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(req)
}

func (h *Handler) delete(ctx *gin.Context) ginx.Render {
	tenantID := pkg.TenantIDFromGinContext(ctx)
	if tenantID ==  uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}


	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	if err := h.srv.Delete(ctx.Request.Context(), tenantID, id); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

// parseUintParam is a helper function to parse a uint from the URL parameters.
func parseUintParam(ctx *gin.Context, paramName string) (uint, error) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}