package bom

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/bom"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler wraps the BOM service to expose it via HTTP handlers.
type Handler struct {
	service bom.Service
	auth middleware.Service

}

// NewHandler creates a new BOM handler.
func NewHandler(service bom.Service,auth middleware.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RouterRegister registers all the routes for the BOM module.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {

	externalGroup := router.Group("/api/v1/external")
	externalGroup.Use(h.auth.Authenticate())

	externalGroup.Handle(http.MethodPost, "/boms", h.create)
	externalGroup.Handle(http.MethodGet, "/boms/:id", h.getByID)
	externalGroup.Handle(http.MethodGet, "/boms/product/:productID", h.list)
	externalGroup.Handle(http.MethodPut, "/boms/:id", h.update)
	externalGroup.Handle(http.MethodPost, "/boms/:id/release", h.release)
	externalGroup.Handle(http.MethodPost, "/boms/:id/obsolete", h.obsolete)
}

func (h *Handler) create(c *gin.Context) ginx.Render {
	var req bom.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)
	}

	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("tenant id not found"))
	}

	resp, err := h.service.Create(c.Request.Context(), tenantID, &req)
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

	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("tenant id not found"))
	}

	resp, err := h.service.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) list(c *gin.Context) ginx.Render {
	productID, err := uuid.Parse(c.Param("productID"))
	if err != nil {
		return ginx.Error(errors.New("invalid product id"))
	}

	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("tenant id not found"))
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	//offset, limit := ginx.GetPagination(c)

	resp, _, err := h.service.List(c.Request.Context(), tenantID, productID, offset, limit)
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

	var req bom.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)
	}

	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("tenant id not found"))
	}

	resp, err := h.service.Update(c.Request.Context(), tenantID, id, &req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) release(c *gin.Context) ginx.Render {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return ginx.Error(err)
	}

	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("tenant id not found"))
	}

	resp, err := h.service.Release(c.Request.Context(), tenantID, id)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) obsolete(c *gin.Context) ginx.Render {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return ginx.Error(err)
	}

	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("tenant id not found"))
	}

	err = h.service.Obsolete(c.Request.Context(), tenantID, id)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

func parseUintParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %s", paramName, idStr)
	}
	return uint(id), nil
}
