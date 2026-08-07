package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API provides gateway-related handlers.
type API struct {
	service *Service
}

// NewAPI creates a new gateway API.
func NewAPI(service *Service) *API {
	return &API{service: service}
}

// RegisterRoutes registers the API routes for gateways.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	gwRoutes := router.Group("/gateways")
	{
		gwRoutes.POST("/register", a.register)
		gwRoutes.GET("", a.list)
		gwRoutes.GET("/:id", a.getByID)
		gwRoutes.POST("/:id/heartbeat", a.heartbeat)
	}
}

// @Summary Register a new gateway
// @Description Creates a new gateway record and returns its details, including the generated ID.
// @Tags gateways
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Gateway Registration Info"
// @Success 201 {object} GatewayResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gateways/register [post]
func (a *API) register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gw, err := a.service.RegisterGateway(c.Request.Context(), req.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ToGatewayResponse(gw))
}

// @Summary List all gateways
// @Description Retrieves a list of all registered gateways.
// @Tags gateways
// @Produce json
// @Success 200 {array} GatewayResponse
// @Failure 500 {object} map[string]string
// @Router /gateways [get]
func (a *API) list(c *gin.Context) {
	gws, err := a.service.ListGateways(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ToGatewayListResponse(gws))
}

// @Summary Get a gateway by ID
// @Description Retrieves details of a specific gateway by its UUID.
// @Tags gateways
// @Produce json
// @Param id path string true "Gateway ID"
// @Success 200 {object} GatewayResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gateways/{id} [get]
func (a *API) getByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gateway ID"})
		return
	}

	gw, err := a.service.GetGateway(c.Request.Context(), id)
	if err != nil {
		// In a real app, you'd check for a "not found" error specifically
		c.JSON(http.StatusNotFound, gin.H{"error": "Gateway not found"})
		return
	}

	c.JSON(http.StatusOK, ToGatewayResponse(gw))
}

// @Summary Gateway heartbeat
// @Description A gateway calls this endpoint periodically to signal it's online.
// @Tags gateways
// @Produce json
// @Param id path string true "Gateway ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gateways/{id}/heartbeat [post]
func (a *API) heartbeat(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gateway ID"})
		return
	}

	if err := a.service.Heartbeat(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}