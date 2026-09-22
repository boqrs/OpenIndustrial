package pkg

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TenantIDFromGinContext(ctx *gin.Context) uuid.UUID {
	value, exists := ctx.Get("tenant_id")
	if !exists {
		return uuid.Nil
	}

	switch value := value.(type) {
	case uuid.UUID:
		return value

	case string:
		id, err := uuid.Parse(value)
		if err != nil {
			return uuid.Nil
		}

		return id

	default:
		return uuid.Nil
	}
}

func GetUserIDFromGinContext(ctx *gin.Context) uuid.UUID {
	value, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil
	}

	switch value := value.(type) {
	case uuid.UUID:
		return value

	case string:
		id, err := uuid.Parse(value)
		if err != nil {
			return uuid.Nil
		}

		return id

	default:
		return uuid.Nil
	}
}

func TenantIDFromContext(ctx context.Context) uuid.UUID {
	value := ctx.Value("tenant_id")

	switch value := value.(type) {
	case uuid.UUID:
		return value

	case string:
		id, err := uuid.Parse(value)
		if err != nil {
			return uuid.Nil
		}

		return id

	default:
		return uuid.Nil
	}
}

func UserIDFromContext(ctx context.Context) uuid.UUID {
	value := ctx.Value("user_id")

	switch value := value.(type) {
	case uuid.UUID:
		return value

	case string:
		id, err := uuid.Parse(value)
		if err != nil {
			return uuid.Nil
		}

		return id

	default:
		return uuid.Nil
	}
}

type BasePageReq struct {
	CurrentPage int `form:"currentPage" json:"currentPage"`
	PageSize    int `form:"pageSize" json:"pageSize"`
}

type PageBaseResp struct {
	Total int64 `json:"total"`
	Next  bool  `json:"next"`
}
