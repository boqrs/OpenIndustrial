package pkg

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

)


func TenantIDFromContext(ctx *gin.Context) uuid.UUID {
	// In a real application, you would extract this from a JWT token or similar.
	val := ctx.Value("tenant_id")
	if val != nil {
		if id, ok := val.(uuid.UUID); ok {
			return id
		}
	}
	// Fallback for testing or unauthenticated contexts
	return uuid.Nil
}


func GetUserIDFromContext(ctx *gin.Context) uuid.UUID {
	// In a real application, you would extract this from a JWT token or similar.
	val := ctx.Value("user_id")
	if val != nil {
		if id, ok := val.(uuid.UUID); ok {
			return id
		}
	}
	// Fallback for testing or unauthenticated contexts
	return uuid.Nil
}
