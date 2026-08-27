package workorder

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"
	//"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	//"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
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
		externalGroup.Handle(http.MethodPost, "/work-orders", h.create)
		externalGroup.Handle(http.MethodGet, "/work-orders", h.list)
		externalGroup.Handle(http.MethodGet, "/work-orders/:id", h.getByID)
		externalGroup.Handle(http.MethodPut, "/work-orders/:id", h.update)
		externalGroup.Handle(http.MethodPost, "/work-orders/:id/release", h.release)
		externalGroup.Handle(http.MethodPost, "/work-orders/:id/start", h.start)
		externalGroup.Handle(http.MethodPost, "/work-orders/:id/cancel", h.cancel)
	}
}


func (h *Handler) create(c *gin.Context) ginx.Render {
	var req srv.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)
	}

	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
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
		return ginx.Error(errors.New("no perm"))
	}

	resp, err := h.service.GetByID(c.Request.Context(), tenantID, id)
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

	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}
	req.TenantID = tenantID

	// Assuming ProductID is a query param for filtering
	productIDStr := c.Query("product_id")
	if productIDStr != "" {
		productID, err := strconv.ParseUint(productIDStr, 10 , 74)
		if err != nil {
			return ginx.Error(errors.New("invalid product_id format"))
		}
		req.ProductID = uint(productID)
	}

	resp, _, err := h.service.List(c.Request.Context(), &req)
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

	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
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
		return ginx.Error(errors.New("no perm"))
	}
	err = h.service.Release(c.Request.Context(), tenantID, id)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

func (h *Handler) start(c *gin.Context) ginx.Render {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return ginx.Error(err)
	}
	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}
	err = h.service.Start(c.Request.Context(), tenantID, id)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

func (h *Handler) complete(c *gin.Context) ginx.Render {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return ginx.Error(err)
	}
	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}
	err = h.service.Complete(c.Request.Context(), tenantID, id)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

func (h *Handler) cancel(c *gin.Context) ginx.Render {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return ginx.Error(err)
	}
	tenantID := pkg.TenantIDFromGinContext(c)
	if tenantID == uuid.Nil {
		return ginx.Error(errors.New("no perm"))
	}
	err = h.service.Cancel(c.Request.Context(), tenantID, id)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

// --- Helpers ---
func parseUintParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %s", paramName, idStr)
	}
	return uint(id), nil
}

