package manufacturing

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/application"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service srv.Service
	auth    middleware.Service
}

func NewHandler(service srv.Service, auth middleware.Service) *Handler {
	return &Handler{service: service, auth: auth}
}

func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {
	group := router.Group("/api/v1/external")
	group.Use(h.auth.Authenticate())

	// --------------------------------------------------------------------- // Production Execution // ---------------------------------------------------------------------
	group.Handle(http.MethodPost, "/executions/execution/:work_order_id", h.CreateProductionExecution)
	group.Handle(http.MethodPost, "/executions/start/:id", h.StartProductionExecution)
	// --------------------------------------------------------------------- // Production Operation // ---------------------------------------------------------------------
	group.Handle(http.MethodPost, "/executions/:execution_id/operations/start/:op_id", h.StartProductionOperation)
	group.Handle(http.MethodPost, "/executions/:execution_id/operations/complete/:op_id", h.CompleteProductionOperation)
	group.Handle(http.MethodPost, "/executions/:execution_id/operations/fail/:op_id", h.FailProductionOperation)
	// --------------------------------------------------------------------- // Production Result // ---------------------------------------------------------------------
	group.Handle(http.MethodPost, "/executions/results/:id/confirm", h.ConfirmExecutionResult)
	group.Handle(http.MethodGet, "/executions", h.ListProductionExecutions)
	group.Handle(http.MethodGet, "/executions/:execution_id", h.GetProductionExecution)
	group.Handle(http.MethodGet, "/executions/:execution_id/operations", h.ListProductionOperations)
}

func (h *Handler) CreateProductionExecution(ctx *gin.Context) ginx.Render {
	workOrderID, err := parseUintParam(ctx, "work_order_id")
	if err != nil {
		return ginx.Error(err)
	}
	result, err := h.service.CreateProductionExecution(ctx.Request.Context(), workOrderID)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(result)
}

func (h *Handler) StartProductionExecution(ctx *gin.Context) ginx.Render {
	executionID, err := parseUintParam(ctx, "execution_id")
	if err != nil {
		return ginx.Error(err)
	}
	if err := h.service.StartProductionExecution(ctx.Request.Context(), executionID); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

func (h *Handler) StartProductionOperation(ctx *gin.Context) ginx.Render {
	executionID, err := parseUintParam(ctx, "execution_id")
	if err != nil {
		return ginx.Error(err)
	}
	operationID, err := parseUintParam(ctx, "op_id")
	if err != nil {
		return ginx.Error(err)
	}
	if err := h.service.StartProductionOperation(ctx.Request.Context(), executionID, operationID); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

type operationResultRequest struct {
	Result map[string]any `json:"result"`
}

func (h *Handler) CompleteProductionOperation(ctx *gin.Context) ginx.Render {
	executionID, err := parseUintParam(ctx, "execution_id")
	if err != nil {
		return ginx.Error(err)
	}
	operationID, err := parseUintParam(ctx, "op_id")
	if err != nil {
		return ginx.Error(err)
	}
	var req operationResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { // Some operations may not produce a result.
		return ginx.Error(err)
	}
	if err := h.service.CompleteProductionOperation(ctx.Request.Context(), executionID, operationID, req.Result); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

func (h *Handler) FailProductionOperation(ctx *gin.Context) ginx.Render {
	executionID, err := parseUintParam(ctx, "execution_id")
	if err != nil {
		return ginx.Error(err)
	}
	operationID, err := parseUintParam(ctx, "op_id")
	if err != nil {
		return ginx.Error(err)
	}
	var req operationResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)

	}
	if err := h.service.FailProductionOperation(ctx.Request.Context(), executionID, operationID, req.Result); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

func (h *Handler) ConfirmExecutionResult(ctx *gin.Context) ginx.Render {
	executionResultID, err := parseUintParam(ctx, "execution_id")
	if err != nil {
		return ginx.Error(err)
	}
	if err := h.service.ConfirmExecutionResult(ctx.Request.Context(), executionResultID); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

func parseUintParam(ctx *gin.Context, paramName string) (uint, error) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: '%s'", paramName, idStr)
	}
	return uint(id), nil
}

// ListProductionExecutions lists executions with optional filters.
func (h *Handler) ListProductionExecutions(ctx *gin.Context) ginx.Render {
	var workOrderID *uint
	if raw := ctx.Query("work_order_id"); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 32)
		if err != nil || id == 0 {
			return ginx.Error(fmt.Errorf("invalid work_order_id: %s", raw))
		}

		parsed := uint(id)
		workOrderID = &parsed
	}

	var status *model.ProductionExecutionStatus
	if raw := ctx.Query("status"); raw != "" {
		parsed := model.ProductionExecutionStatus(raw)
		status = &parsed
	}

	result, err := h.service.ListProductionExecutions(
		ctx.Request.Context(),
		workOrderID,
		status,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// GetProductionExecution retrieves an execution by ID.
func (h *Handler) GetProductionExecution(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "execution_id")
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.GetProductionExecution(
		ctx.Request.Context(),
		id,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

// ListProductionOperations lists operations belonging to an execution.
func (h *Handler) ListProductionOperations(ctx *gin.Context) ginx.Render {
	id, err := parseUintParam(ctx, "execution_id")
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.ListProductionOperations(
		ctx.Request.Context(),
		id,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}
