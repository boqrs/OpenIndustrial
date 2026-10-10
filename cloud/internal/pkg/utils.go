package pkg

import (
// "context"

// "github.com/gin-gonic/gin"
// "github.com/google/uuid"
)

// type tenantIDContextKey struct{}

// // WithTenantID stores the authenticated tenant ID in a standard context.
// func WithTenantID(ctx context.Context, tenantID uint) context.Context {
// 	return context.WithValue(ctx, tenantIDContextKey{}, tenantID)
// }

// // TenantIDUintFromContext returns the authenticated tenant's internal ID.
// func TenantIDUintFromContext(ctx context.Context) uint {
// 	value := ctx.Value(tenantIDContextKey{})

// 	tenantID, ok := value.(uint)
// 	if !ok {
// 		return 0
// 	}

// 	return tenantID
// }

// func TenantIDUintFromGinContext(ctx *gin.Context) uint {
// 	value := ctx.Value("tenant_id")
// 	tenantID, ok := value.(uint)
// 	if !ok {
// 		return 0
// 	}

// 	return tenantID
// }

// func GetUserIDFromGinContext(ctx *gin.Context) uuid.UUID {
// 	value, exists := ctx.Get("user_id")
// 	if !exists {
// 		return uuid.Nil
// 	}

// 	switch value := value.(type) {
// 	case uuid.UUID:
// 		return value

// 	case string:
// 		id, err := uuid.Parse(value)
// 		if err != nil {
// 			return uuid.Nil
// 		}

// 		return id

// 	default:
// 		return uuid.Nil
// 	}
// }

// func UserIDFromContext(ctx context.Context) uuid.UUID {
// 	value := ctx.Value("user_id")

// 	switch value := value.(type) {
// 	case uuid.UUID:
// 		return value

// 	case string:
// 		id, err := uuid.Parse(value)
// 		if err != nil {
// 			return uuid.Nil
// 		}

// 		return id

// 	default:
// 		return uuid.Nil
// 	}
// }

type BasePageReq struct {
	CurrentPage int `form:"currentPage" json:"currentPage"`
	PageSize    int `form:"pageSize" json:"pageSize"`
}

type PageBaseResp struct {
	Total int64 `json:"total"`
	Next  bool  `json:"next"`
}
