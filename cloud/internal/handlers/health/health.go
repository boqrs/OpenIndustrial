package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHealthRoutes registers a simple /health endpoint.
// This is useful for load balancers and monitoring systems to check if the service is alive.
func RegisterHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}