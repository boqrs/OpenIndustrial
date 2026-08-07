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
	// For now, we'll use a placeholder OrgID. In a real app, this would come from auth middleware.
	orgID := uuid.New()

	var req CreateWorkOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wo, err := a.service.CreateWorkOrder(c.Request.Context(), orgID, req.ProductID, req.Quantity, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create work order"})
		return
	}

	c.JSON(http.StatusCreated, ToWorkOrderResponse(wo))
}

func (a *API) listWorkOrders(c *gin.Context) {
	// Placeholder OrgID
	orgID := uuid.New()

	wos, err := a.service.ListByOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list work orders"})
		return
	}

	c.JSON(http.StatusOK, ToWorkOrderListResponse(wos))
}

func (a *API) getWorkOrder(c *gin.Context) {
	// Placeholder OrgID
	orgID := uuid.New()
	workOrderID, err := uuid.Parse(c.Param("workorder_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work order ID"})
		return
	}

	wo, err := a.service.FindByID(c.Request.Context(), orgID, workOrderID)
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
	// Placeholder OrgID
	orgID := uuid.New()
	workOrderID, err := uuid.Parse(c.Param("workorder_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work order ID"})
		return
	}

	if err := a.service.StartWorkOrder(c.Request.Context(), orgID, workOrderID); err != nil {
		if err == ErrWorkOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start work order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Work order started"})
}