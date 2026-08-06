package workorder

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API provides the HTTP handlers for the work order domain.
type API struct {
	service *Service
}

// NewAPI creates a new work order API handler.
func NewAPI(service *Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the work order API routes with a Gin router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	woRoutes := router.Group("/workorders")
	{
		woRoutes.POST("", a.createWorkOrder)
		woRoutes.GET("", a.listWorkOrders)
		woRoutes.GET("/:workorder_id", a.getWorkOrder)
		woRoutes.POST("/:workorder_id/start", a.startWorkOrder)
	}
}

func (a *API) createWorkOrder(c *gin.Context) {
	orgIDStr := c.Param("org_id") // Assuming org_id is a URL parameter
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var req CreateWorkOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	wo, err := a.service.CreateWorkOrder(c.Request.Context(), orgID, productID, req.Quantity, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create work order"})
		return
	}

	c.JSON(http.StatusCreated, ToWorkOrderResponse(wo))
}

func (a *API) listWorkOrders(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	wos, err := a.service.repo.ListByOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list work orders"})
		return
	}

	c.JSON(http.StatusOK, ToWorkOrderListResponse(wos))
}

func (a *API) getWorkOrder(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	workOrderIDStr := c.Param("workorder_id")
	workOrderID, err := uuid.Parse(workOrderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work order id"})
		return
	}

	wo, err := a.service.repo.FindByID(c.Request.Context(), orgID, workOrderID)
	if err != nil {
		if err == ErrWorkOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get work order"})
		return
	}

	c.JSON(http.StatusOK, ToWorkOrderResponse(wo))
}

func (a *API) startWorkOrder(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	workOrderIDStr := c.Param("workorder_id")
	workOrderID, err := uuid.Parse(workOrderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work order id"})
		return
	}

	wo, err := a.service.StartWorkOrder(c.Request.Context(), orgID, workOrderID)
	if err != nil {
		// Handle specific errors like "already started" in the future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start work order"})
		return
	}

	c.JSON(http.StatusOK, ToWorkOrderResponse(wo))
}