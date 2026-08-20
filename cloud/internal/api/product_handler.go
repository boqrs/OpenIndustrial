package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/OpenIndustrial/cloud/internal/product"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API handles HTTP requests for the product module.
type ProductAPI struct {
	service product.Service
}

// NewAPI creates a new API handler for the product service.
func NewProductAPI(service product.Service) *ProductAPI {
	return &ProductAPI{service: service}
}

// Register registers all product model routes to the given router group.
func (a *ProductAPI) Register(group *gin.RouterGroup) {
	productModels := group.Group("/product-models")
	{
		productModels.POST("", a.createProductModel)
		productModels.GET("", a.listProductModels)
		productModels.GET("/:id", a.getProductModel)
		productModels.PATCH("/:id", a.updateProductModel)

		// Attribute management
		productModels.GET("/:id/attributes", a.getAttributes)
		productModels.PUT("/:id/attributes", a.updateAttributes)

		// Status management
		productModels.POST("/:id/activate", a.activateProductModel)
		productModels.POST("/:id/deactivate", a.deactivateProductModel)
		productModels.POST("/:id/archive", a.archiveProductModel)
	}
}

func (a *ProductAPI) createProductModel(c *gin.Context) {
	var req param.CreateProductModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	resp, err := a.service.CreateProductModel(c.Request.Context(), &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (a *ProductAPI) listProductModels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	req := param.ListProductModelsRequest{
		Category: c.Query("category"),
		Status:   c.Query("status"),
		Code:     c.Query("code"),
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := a.service.ListProductModels(c.Request.Context(), &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (a *ProductAPI) getProductModel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product model ID format"})
		return
	}

	resp, err := a.service.GetProductModel(c.Request.Context(), id)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (a *ProductAPI) updateProductModel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product model ID format"})
		return
	}

	var req param.UpdateProductModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	resp, err := a.service.UpdateProductModel(c.Request.Context(), id, &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (a *ProductAPI) getAttributes(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product model ID format"})
		return
	}

	resp, err := a.service.GetAttributeDefinitions(c.Request.Context(), id)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (a *ProductAPI) updateAttributes(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product model ID format"})
		return
	}

	var req param.UpdateAttributeDefinitionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	err = a.service.UpdateAttributeDefinitions(c.Request.Context(), id, &req)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (a *ProductAPI) activateProductModel(c *gin.Context) {
	a.updateStatus(c, model.StatusActive)
}

func (a *ProductAPI) deactivateProductModel(c *gin.Context) {
	a.updateStatus(c, model.StatusInactive)
}

func (a *ProductAPI) archiveProductModel(c *gin.Context) {
	a.updateStatus(c, model.StatusArchived)
}

func (a *ProductAPI) updateStatus(c *gin.Context, status string) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product model ID format"})
		return
	}

	err = a.service.UpdateProductModelStatus(c.Request.Context(), id, status)
	if err != nil {
		a.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleError centralizes error handling for the product API.
func (a *ProductAPI) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, product.ErrProductModelNotFound), errors.Is(err, product.ErrAttributeDefinitionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, product.ErrProductModelCodeVersionExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, product.ErrProductModelImmutable),
		errors.Is(err, product.ErrInvalidProductModel),
		errors.Is(err, product.ErrInvalidAttributeDefinition),
		errors.Is(err, product.ErrProductModelAlreadyActive),
		errors.Is(err, product.ErrProductModelCannotModify):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		// Log the error for internal review
		// log.Printf("internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
	}
}