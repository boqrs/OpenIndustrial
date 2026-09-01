package pkg

import (
	"context"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

)


func TenantIDFromGinContext(ctx *gin.Context) uuid.UUID {
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


func TenantIDFromContext(ctx context.Context) uuid.UUID {
	value := ctx.Value("tenant_id")

	if id, ok := value.(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}

type BasePageReq struct {
	CurrentPage int `form:"currentPage" json:"currentPage"`
	PageSize    int `form:"pageSize" json:"pageSize"`
}

type PageBaseResp struct {
	Total int64 `json:"total"`
	Next  bool  `json:"next"`
}