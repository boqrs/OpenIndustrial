package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// NewAuthMiddleware creates a middleware for JWT authentication.
func NewAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Set user and tenant info in the context for downstream handlers
		c.Set("user_id", claims["sub"])
		c.Set("tenant_id", claims["tid"])

		c.Next()
	}
}

// --- Context Helper Functions ---

func getTenantIDFromContext(c *gin.Context) (uuid.UUID, error) {
	tenantIDStr, exists := c.Get("tenant_id")
	if !exists {
		return uuid.Nil, errors.New("tenant_id not found in context")
	}
	tenantID, err := uuid.Parse(tenantIDStr.(string))
	if err != nil {
		return uuid.Nil, errors.New("invalid tenant_id format in context")
	}
	return tenantID, nil
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user_id not found in context")
	}
	return uuid.Parse(userIDStr.(string))
}