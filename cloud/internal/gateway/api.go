package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API provides the HTTP handlers for the gateway domain.
type API struct {
	service *Service
}

// NewAPI creates a new gateway API handler.
func NewAPI(service *Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the gateway API routes with a Gin router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	gatewayRoutes := router.Group("/gateways")
	{
		// This endpoint is for gateways to register themselves
		gatewayRoutes.POST("/register", a.registerGateway)
		// This is for gateways to send heartbeats
		gatewayRoutes.POST("/:gateway_id/heartbeat", a.handleHeartbeat)

		// These are for users to manage gateways
		gatewayRoutes.GET("", a.listGateways)
	}
}

func (a *API) registerGateway(c *gin.Context) {
	orgIDStr := c.GetHeader("X-Org-ID") // Gateway auth might use headers
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id in header"})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	gatewayID, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gateway id"})
		return
	}

	gw, err := a.service.RegisterGateway(c.Request.Context(), orgID, gatewayID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register gateway"})
		return
	}

	c.JSON(http.StatusOK, ToGatewayResponse(gw))
}

func (a *API) handleHeartbeat(c *gin.Context) {
	orgIDStr := c.GetHeader("X-Org-ID")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id in header"})
		return
	}

	gatewayIDStr := c.Param("gateway_id")
	gatewayID, err := uuid.Parse(gatewayIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gateway id"})
		return
	}

	if err := a.service.HandleHeartbeat(c.Request.Context(), orgID, gatewayID); err != nil {
		if err == ErrGatewayNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to handle heartbeat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *API) listGateways(c *gin.Context) {
	// This would be a user-facing endpoint, so it uses user's org from context
	orgIDStr := c.Param("org_id") // Or from JWT context: c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	gws, err := a.service.repo.ListByOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list gateways"})
		return
	}

	c.JSON(http.StatusOK, ToGatewayListResponse(gws))
}