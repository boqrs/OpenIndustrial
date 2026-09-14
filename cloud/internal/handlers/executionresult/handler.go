package executionresult

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/executionresult"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
)

// Handler manages execution result endpoints.
type Handler struct {
	service srv.Service
	auth    middleware.Service
}

// NewHandler creates a new execution result handler.
func NewHandler(
	service srv.Service,
	auth middleware.Service,
) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
	}
}

// RouterRegister registers execution result routes.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {
	group := router.Group("/api/v1/external")
	group.Use(h.auth.Authenticate())

	group.Handle(
		http.MethodPost,
		"/execution-results",
		h.CreateResult,
	)

	group.Handle(
		http.MethodGet,
		"/execution-results/:id",
		h.GetResult,
	)

	group.Handle(
		http.MethodGet,
		"/executions/:execution_id/result",
		h.GetResultByExecutionID,
	)

	group.Handle(
		http.MethodPost,
		"/execution-results/:id/cancel",
		h.CancelResult,
	)
}

// CreateResult creates a production execution result.
func (h *Handler) CreateResult(ctx *gin.Context) ginx.Render {
	var req srv.CreateResultRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(
			fmt.Errorf("invalid request payload: %w", err),
		)
	}

	result, err := h.service.CreateResult(
		ctx.Request.Context(),
		&req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// GetResult retrieves an execution result by ID.
func (h *Handler) GetResult(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.GetResult(
		ctx.Request.Context(),
		id,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// GetResultByExecutionID retrieves the result for an execution.
func (h *Handler) GetResultByExecutionID(ctx *gin.Context) ginx.Render {
	executionID, err := parseUintParam(ctx, "execution_id")
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.GetResultByExecutionID(
		ctx.Request.Context(),
		executionID,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// CancelResult cancels a draft execution result.
func (h *Handler) CancelResult(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	if err := h.service.CancelResult(
		ctx.Request.Context(),
		id,
	); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

func parseUintParam(ctx *gin.Context, paramName string) (uint, error) {
	idStr := ctx.Param(paramName)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf(
			"invalid %s format: '%s'",
			paramName,
			idStr,
		)
	}

	return uint(id), nil
}
