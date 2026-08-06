package auth

import (
	"net/http"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/user"
	"github.com/gin-gonic/gin"
)

// API provides the HTTP handlers for the auth domain.
type API struct {
	service *Service
}

// NewAPI creates a new auth API handler.
func NewAPI(service *Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the auth API routes with a Gin router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", a.login)
	}
}

func (a *API) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.service.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == user.ErrUserNotFound || err == ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}