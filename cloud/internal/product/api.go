package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API provides the product-related endpoints.
type API struct {
	service *Service
}

// NewAPI creates a new product API handler.
func NewAPI(service *Service) *API {
	return &API{service: service}
}

// RegisterRoutes registers the product routes on the given router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	productRoutes := router.Group("/products")
	{
		productRoutes.POST("", a.createProduct)
		productRoutes.GET("", a.listProducts)
		productRoutes.GET("/:id", a.getProduct)
	}
}

// createProduct handles the creation of a new product.
func (a *API) createProduct(c *gin.Context) {
	// For now, we'll use a placeholder OrgID. In a real app, this would come from auth middleware.
	orgID := uuid.New()

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := a.service.CreateProduct(c.Request.Context(), orgID, req.Name, req.Code, req.Spec, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// listProducts handles listing products for an organization.
func (a *API) listProducts(c *gin.Context) {
	// Placeholder OrgID
	orgID := uuid.New()

	products, err := a.service.ListProductsByOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// getProduct handles retrieving a single product by its ID.
func (a *API) getProduct(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	product, err := a.service.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		// This could be a 404 Not Found in a real system
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}

	c.JSON(http.StatusOK, product)
}