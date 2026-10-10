package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/product"
	zlog "github.com/boqrs/nexus/log"
	"github.com/boqrs/zeus/ginx"
)

// API handles HTTP requests for the product module.
type Handler struct {
	service srv.Service
	l *zlog.Provider
}

// NewAPI creates a new API handler for the product service.
func NewHandler(service srv.Service, logger *zlog.Provider) *Handler {
	return &Handler{
		service: service,
		l:logger,
	}
}

// Register registers all product model routes to the given router group.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {
	externalGroup := router.Group("/api/v1/external")
	externalGroup.Handle(http.MethodPost, "/product-models", h.createProductModel)
	externalGroup.Handle(http.MethodGet, "/product-models", h.listProductModels)
	externalGroup.Handle(http.MethodGet, "/product-models/:id", h.getProductModel)
	externalGroup.Handle(http.MethodPatch, "/product-models/:id", h.updateProductModel)
	externalGroup.Handle(http.MethodGet, "/product-models/:id/attributes", h.getAttributes)
	externalGroup.Handle(http.MethodPut, "/product-models/:id/attributes", h.updateAttributes)
	externalGroup.Handle(http.MethodPost, "/product-models/:id/activate", h.activateProductModel)
	externalGroup.Handle(http.MethodPost, "/product-models/:id/deactivate", h.deactivateProductModel)
	externalGroup.Handle(http.MethodPost, "/product-models/:id/archive", h.archiveProductModel)
}

func (h *Handler) createProductModel(ctx *gin.Context) ginx.Render {
	var req srv.CreateProductModelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := h.service.CreateProductModel(ctx.Request.Context(), &req)
	if err != nil {
		h.l.Get().Errorf("failed to create product model, error: %s", err.Error())
		return ginx.Error(err)
	}

	return ginx.Success(resp)

}

func (h *Handler) listProductModels(ctx *gin.Context) ginx.Render {
	req := srv.ListProductModelsRequest{}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := h.service.ListProductModels(ctx.Request.Context(), &req)
	if err != nil {
		h.l.Get().Errorf("failed to get product list, error: %s", err.Error())
		return ginx.Error(err)
	}

	return ginx.Success(resp)

}

func (h *Handler) getProductModel(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := h.service.GetProductModel(ctx.Request.Context(), id)
	if err != nil {
		h.l.Get().Errorf("failed to get product, error: %s", err.Error())
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) updateProductModel(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	var req srv.UpdateProductModelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}

	resp, err := h.service.UpdateProductModel(ctx.Request.Context(), id, &req)
	if err != nil {
		h.l.Get().Errorf("failed to update product, error: %s", err.Error())
		return ginx.Error(err)
	}

	return ginx.Success(resp)

}

func (h *Handler) getAttributes(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	resp, err := h.service.GetAttributeDefinitions(ctx.Request.Context(), id)
	if err != nil {
		h.l.Get().Errorf("failed to get product attributes, error: %s", err.Error())
		return ginx.Error(err)
	}

	return ginx.Success(resp)

}

func (h *Handler) updateAttributes(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid param"))
	}

	var req srv.UpdateAttributeDefinitionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}

	if err = h.service.UpdateAttributeDefinitions(ctx.Request.Context(), id, &req);err != nil {
		h.l.Get().Errorf("failed to update product attributes, error: %s", err.Error())
		return ginx.Error(err)
	}

	return ginx.Success(nil)

}

func (h *Handler) activateProductModel(ctx *gin.Context) ginx.Render {
	h.updateStatus(ctx, model.StatusActive)
	return ginx.Success(nil)
}

func (h *Handler) deactivateProductModel(ctx *gin.Context) ginx.Render {
	h.updateStatus(ctx, model.StatusInactive)
	return ginx.Success(nil)
}

func (h *Handler) archiveProductModel(ctx *gin.Context) ginx.Render {
	h.updateStatus(ctx, model.StatusArchived)
	return ginx.Success(nil)
}

func (h *Handler) updateStatus(ctx *gin.Context, status string) {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return
	}

	if err = h.service.UpdateProductModelStatus(ctx.Request.Context(), id, status); err != nil {
		h.l.Get().Errorf("failed to update status: %s, error: %s", status, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Helper function to parse uint from path parameter
func parseUintParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %s", paramName, idStr)
	}
	return uint(id), nil
}
