package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", a.CreateUser)
	}
}

// CreateUser handles the HTTP request to create a new user.
func (a *API) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.service.RegisterNewUser(c.Request.Context(), req.OrgID, req.Username, req.Password, req.Email, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ToUserResponse(user))
}