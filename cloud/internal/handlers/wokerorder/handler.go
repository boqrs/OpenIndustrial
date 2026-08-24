package workorder

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	"github.com/boqrs/zeus/ginx"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Handler struct {
	service srv.Service
	auth middleware.Service
}

func NewHandler(service srv.Service, auth middleware.Service) *Handler {
	return &Handler{
		service: service,
		auth:auth,
	}
}

func (h *Handler)  RouterRegister(router ginx.ZeroGinRouter){
	// External APIs
	externalGroup := router.Group("/api/v1/external")
	//externalGroup.Use(mid.CORS())
	// 首选需要 auth
	externalGroup.Use(h.auth.Authenticate())
		externalGroup.Handle(http.MethodPost, "/work-orders", h.Create)
		externalGroup.Handle(http.MethodGet, "/work-orders", h.List)
		externalGroup.Handle(http.MethodGet, "/work-orders/:id", h.Get)
		externalGroup.Handle(http.MethodPut, "/work-orders/:id", h.Update)
		externalGroup.Handle(http.MethodPost, "/work-orders/:id/release", h.Release)
		externalGroup.Handle(http.MethodPost, "/work-orders/:id/start", h.Start)
		externalGroup.Handle(http.MethodPost, "/work-orders/:id/cancel", h.Cancel)
}



func (h *Handler) Create(ctx *gin.Context) ginx.Render{
	var req srv.CreateWorkOrderRequest
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil{
		return ginx.Error(fmt.Errorf("no perm"))
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid parm json"))
	}

	result, err := h.service.CreateWorkOrder(ctx.Request.Context(),&req)
	if err != nil {
		//handleError(c, err)
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

func (h *Handler) Get(ctx *gin.Context) ginx.Render {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.GetWorkOrder(ctx.Request.Context(),id)
	if err != nil {
		//handleError(c, err)
		return ginx.Error(err)
	}
	return ginx.Success(result)
}

func (h *Handler) List(ctx *gin.Context) ginx.Render {
	var status *model.WorkOrderStatus
	if value := ctx.Query("status"); value != "" {
		s := model.WorkOrderStatus(value)
		status = &s
	}

	var productionPlanID *uuid.UUID
	if value := ctx.Query("production_plan_id"); value != "" {
		id, err := uuid.Parse(value)

		if err != nil {
			return ginx.Error(err)
		}

		productionPlanID = &id
	}

	result, err := h.service.ListWorkOrders(
		ctx.Request.Context(),
		status,
		productionPlanID,
	)

	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(result)
}

func (h *Handler) Update(ctx *gin.Context) ginx.Render {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ginx.Error(err)
	}

	var req srv.UpdateWorkOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(err)
	}

	result, err := h.service.UpdateWorkOrder(ctx.Request.Context(),id,&req)
	if err != nil {
		//handleError(c, err)
		return ginx.Error(err)
	}
	return ginx.Success(result)
}

func (h *Handler) Release(ctx *gin.Context) ginx.Render {
	h.changeState(ctx,
		func(id uuid.UUID) error {
			return h.service.ReleaseWorkOrder(
				ctx.Request.Context(),
				id,
			)
		},
	)

	return ginx.Success(nil)
}

func (h *Handler) Start(ctx *gin.Context) ginx.Render {
	h.changeState(
		ctx,
		func(id uuid.UUID) error {
			return h.service.StartWorkOrder(
				ctx.Request.Context(),
				id,
			)
		},
	)
	return ginx.Success(nil)
}

func (h *Handler) Cancel(ctx *gin.Context) ginx.Render {
	h.changeState(
		ctx,
		func(id uuid.UUID) error {
			return h.service.CancelWorkOrder(
				ctx.Request.Context(),
				id,
			)
		},
	)

	return ginx.Success(nil)
}

func (h *Handler) changeState(c *gin.Context,fn func(uuid.UUID) error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return 
	}

	if err := fn(id); err != nil {
		//handleError(c, err)
		return
	}

	//c.Status(http.StatusNoContent)
}