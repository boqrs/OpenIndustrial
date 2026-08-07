package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// API provides the HTTP handlers for the product domain.
type API struct {
	service Service
}

// NewAPI creates a new product API handler.
func NewAPI(service Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the product API routes with a Gin router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	productRoutes := router.Group("/products")
	{
		productRoutes.POST("", a.createProduct)
		productRoutes.GET("", a.listProducts)
		productRoutes.GET("/:product_id", a.getProduct)
	}
}

func (a *API) createProduct(c *gin.Context) {
	orgIDStr := c.Param("org_id") // Assuming org_id is a URL parameter

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := a.service.CreateProduct(c.Request.Context(), orgIDStr, req.Name, req.Description, req.Spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, ToProductResponse(p))
}

func (a *API) listProducts(c *gin.Context) {
	orgIDStr := c.Param("org_id")

	products, err := a.service.ListProductsForOrg(c.Request.Context(), orgIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}

	c.JSON(http.StatusOK, ToProductListResponse(products))
}

func (a *API) getProduct(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	productIDStr := c.Param("product_id")

	p, err := a.service.GetProductByID(c.Request.Context(), orgIDStr, productIDStr)
	if err != nil {
		// Consider adding a specific ErrProductNotFound
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, ToProductResponse(p))
}