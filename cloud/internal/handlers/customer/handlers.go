package customer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/customer"
	"github.com/boqrs/zeus/ginx"
)

type Handler struct {
	service srv.Service
	auth    middleware.Service
}

func NewHandler(
	service srv.Service,
	auth middleware.Service,
) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
	}
}

func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {

	externalGroup := router.Group("/api/v1/external")
	externalGroup.Use(h.auth.Authenticate())
	externalGroup.Handle(http.MethodPost, "/customers", h.create)
	externalGroup.Handle(http.MethodGet, "/customers", h.list)
	externalGroup.Handle(http.MethodGet, "/customers/code/:code", h.getByCode)
	externalGroup.Handle(http.MethodGet, "/customers/:id", h.getByID)
	externalGroup.Handle(http.MethodPut, "/customers/:id", h.update)
	externalGroup.Handle(http.MethodPost, "/customers/:id/deactivate", h.deactivate)

}

func (h *Handler) create(c *gin.Context) ginx.Render {
	var req srv.CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		return ginx.Error(err)
	}

	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		return ginx.Error(err)

	}

	return ginx.Success(resp)

}

func (h *Handler) list(c *gin.Context) ginx.Render {
	var req srv.ListRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		return ginx.Error(err)
	}

	resp, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		return ginx.Error(err)

	}

	return ginx.Success(resp)

}

func (h *Handler) getByID(c *gin.Context) ginx.Render {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return ginx.Error(err)

	}

	resp, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		return ginx.Error(err)

	}

	return ginx.Success(resp)

}

func (h *Handler) getByCode(c *gin.Context) ginx.Render {
	code := c.Param("code")

	resp, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		return ginx.Error(err)

	}

	return ginx.Success(resp)

}

func (h *Handler) update(c *gin.Context) ginx.Render {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return ginx.Error(err)

	}

	var req srv.UpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)

	}

	resp, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		return ginx.Error(err)

	}

	return ginx.Success(resp)

}

func (h *Handler) deactivate(c *gin.Context) ginx.Render {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return ginx.Error(err)

	}

	if err := h.service.Deactivate(c.Request.Context(), id); err != nil {
		return ginx.Error(err)

	}

	return ginx.Success(nil)

}

func parseUintParam(c *gin.Context, name string) (uint, error) {
	value := c.Param(name)

	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, strconv.ErrSyntax
	}

	return uint(id), nil

}
