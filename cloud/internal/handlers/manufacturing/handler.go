package manufacturing

import (
	"fmt"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/application"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	group.Handle(http.MethodPost, "/work-orders/:work_order_id/execution", h.CreateProductionExecution)
	group.Handle(http.MethodPost, "/executions/:id/start", h.StartProductionExecution)
	// --------------------------------------------------------------------- // Production Operation // ---------------------------------------------------------------------
	group.Handle(http.MethodPost, "/executions/:id/operations/:op_id/start", h.StartProductionOperation)
	group.Handle(http.MethodPost, "/executions/:id/operations/:op_id/complete", h.CompleteProductionOperation)
	group.Handle(http.MethodPost, "/executions/:id/operations/:op_id/fail", h.FailProductionOperation)
	// --------------------------------------------------------------------- // Production Result // ---------------------------------------------------------------------
	group.Handle(http.MethodPost, "/execution-results/:id/confirm", h.ConfirmExecutionResult)
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
	executionID, err := parseUintParam(ctx, "id")
	if err != nil {
		return ginx.Error(err)
	}
	if err := h.service.StartProductionExecution(ctx.Request.Context(), executionID); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

func (h *Handler) StartProductionOperation(ctx *gin.Context) ginx.Render {
	executionID, err := parseUintParam(ctx, "id")
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
	executionID, err := parseUintParam(ctx, "id")
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
	executionID, err := parseUintParam(ctx, "id")
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
	executionResultID, err := parseUintParam(ctx, "id")
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
