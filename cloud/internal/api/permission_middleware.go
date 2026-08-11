package api

import (
	"net/http"

	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PermissionRepository is now an alias for the definitive interface in the identity package.
// This ensures consistency and that our middleware uses the correct, implemented methods.
type PermissionRepository interface {
	identity.PermissionRepository
}

// PermissionMiddleware is a struct that holds the dependency for checking permissions.
type PermissionMiddleware struct {
	repo PermissionRepository
}

// NewPermissionMiddleware creates a new instance of the permission middleware.
func NewPermissionMiddleware(repo PermissionRepository) *PermissionMiddleware {
	return &PermissionMiddleware{repo: repo}
}

// RequirePermission returns a Gin handler function that checks for a specific permission.
// This is the core of our fine-grained authorization.
func (m *PermissionMiddleware) RequirePermission(permissionKey string) gin.HandlerFunc {
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
		hasPermission, err := m.repo.CheckPermissionForUser(c.Request.Context(), userID, permissionKey)
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