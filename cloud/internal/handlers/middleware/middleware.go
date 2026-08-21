package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/identity"
)

type service struct{
	jwtSecret string
	repo identity.PermissionRepository
}

type Service interface{
	Authenticate() gin.HandlerFunc
	RequirePermission(permissionKey string) gin.HandlerFunc
}

func NewAuthService(jwtSecret string, repo identity.PermissionRepository)Service{
	return &service{
		jwtSecret: jwtSecret,
		repo: repo,
	}
}

// NewAuthMiddleware creates a middleware for JWT authentication.
func(s *service) Authenticate() gin.HandlerFunc {
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
			return []byte(s.jwtSecret), nil
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

// RequirePermission returns a Gin handler function that checks for a specific permission.
// This is the core of our fine-grained authorization.
func (s *service) RequirePermission(permissionKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get user ID from the context, which was set by the auth middleware.
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		userID, err := uuid.Parse(userIDRaw.(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format in token"})
			return
		}

		// 2. Check permission using the repository.
		hasPermission, err := s.repo.CheckPermissionForUser(c.Request.Context(), userID, permissionKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			return
		}

		// 3. Enforce permission.
		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have the required permission: " + permissionKey})
			return
		}

		// 4. If permission is granted, proceed to the next handler.
		c.Next()
	}
}

// --- Context Helper Functions ---

func GetTenantIDFromContext(c *gin.Context) (uuid.UUID, error) {
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

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user_id not found in context")
	}
	return uuid.Parse(userIDStr.(string))
}