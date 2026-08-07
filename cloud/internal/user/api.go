package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API encapsulates all the handlers for the user resource.
type API struct {
	service *Service
}

// NewAPI creates a new user API handler.
func NewAPI(service *Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the user API routes to the gin router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("", a.create)
	}
}

// CreateUser handles the HTTP request to create a new user.
func (a *API) create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, err := uuid.Parse(req.OrgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	user, err := a.service.CreateUser(
		c.Request.Context(),
		orgID,
		req.Username,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ToUserResponse(user))
}