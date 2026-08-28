package execution

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
)

// Handler manages the execution endpoints.
type Handler struct {
	service srv.Service
	auth    middleware.Service
}

// NewHandler creates a new execution handler.
func NewHandler(service srv.Service, auth middleware.Service) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
	}
}

// RouterRegister registers the execution routes.
func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {
	execGroup := router.Group("/api/v1/external")
	execGroup.Use(h.auth.Authenticate())
	//{
		execGroup.Handle(http.MethodGet, "/executions", h.ListExecutions)
		execGroup.Handle(http.MethodGet, "/executions/:id", h.GetExecution)
		execGroup.Handle(http.MethodPost, "/executions/:id/start", h.StartExecution)
		execGroup.Handle(http.MethodPost, "/executions/:id/cancel", h.CancelExecution)

		execGroup.Handle(http.MethodGet, "/executions/:id/operations", h.ListOperations)
		execGroup.Handle(http.MethodPost, "/executions/:id/operations/:op_id/start", h.StartOperation)
		execGroup.Handle(http.MethodPost, "/executions/:id/operations/:op_id/complete", h.CompleteOperation)
		execGroup.Handle(http.MethodPost, "/executions/:id/operations/:op_id/fail", h.FailOperation)
	//}
}

// GetExecution retrieves a single production execution.
func (h *Handler) GetExecution(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.GetExecution(ctx.Request.Context(), id)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(result)
}

type listExecutionsRequest struct {
	WorkOrderID *uint                            `form:"workOrderID"`
	DeviceID    *uint                            `form:"deviceID"`
	Status      *model.ProductionExecutionStatus `form:"status"`
}

// ListExecutions retrieves a list of production executions.
func (h *Handler) ListExecutions(ctx *gin.Context) ginx.Render {
	var req listExecutionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid query parameters: %w", err))
	}

	result, err := h.service.ListExecutions(ctx.Request.Context(), req.WorkOrderID, req.DeviceID, req.Status)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(result)
}

// StartExecution starts a production execution.
func (h *Handler) StartExecution(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	if err := h.service.StartExecution(ctx.Request.Context(), id); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

// CancelExecution cancels a production execution.
func (h *Handler) CancelExecution(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	if err := h.service.CancelExecution(ctx.Request.Context(), id); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

// ListOperations retrieves the operations for a specific execution.
func (h *Handler) ListOperations(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.ListOperations(ctx.Request.Context(), id)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(result)
}

// StartOperation starts a specific operation within an execution.
func (h *Handler) StartOperation(ctx *gin.Context) ginx.Render {
	executionID, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}
	operationID, err := parseUintParam(ctx, "op_id")
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid operation ID format: %w", err))
	}

	if err := h.service.StartOperation(ctx.Request.Context(), executionID, operationID); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

type operationResultRequest struct {
	Result map[string]any `json:"result"`
}

// CompleteOperation completes a specific operation with a 'passed' status.
func (h *Handler) CompleteOperation(ctx *gin.Context) ginx.Render {
	executionID, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}
	operationID, err := parseUintParam(ctx, "op_id")
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid operation ID format: %w", err))
	}

	var req operationResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Allow empty body for operations that don't produce a result
		if err.Error() != "EOF" {
			return ginx.Error(fmt.Errorf("invalid request payload: %w", err))
		}
	}

	if err := h.service.CompleteOperation(ctx.Request.Context(), executionID, operationID, req.Result); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

// FailOperation completes a specific operation with a 'failed' status.
func (h *Handler) FailOperation(ctx *gin.Context) ginx.Render {
	executionID, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}
	operationID, err := parseUintParam(ctx, "op_id")
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid operation ID format: %w", err))
	}

	var req operationResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if err.Error() != "EOF" {
			return ginx.Error(fmt.Errorf("invalid request payload: %w", err))
		}
	}

	if err := h.service.FailOperation(ctx.Request.Context(), executionID, operationID, req.Result); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

// parseUintParam is a helper to parse a uint64 from a URL parameter.
func parseUintParam(ctx *gin.Context, paramName string) (uint, error) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: '%s'", paramName, idStr)
	}
	return uint(id), nil
}

