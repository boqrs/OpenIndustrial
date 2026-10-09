package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/identity"

	"github.com/boqrs/zeus/ginx"
)

type Handler struct {
	service srv.Service
	auth    middleware.Service
}

func NewDashboardHandler(
	service srv.Service,
	auth middleware.Service,
) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
	}
}

func (h *Handler) RouterRegister(
	router ginx.ZeroGinRouter,
) {
	adminGroup := router.Group("/api/v1/external")
	adminGroup.Use(
		h.auth.Authenticate(),
		h.auth.RequireAdmin(),
	)

	adminGroup.Handle(
		http.MethodGet,
		"/dashboard/overview",
		h.handleOverview,
	)
}

func (h *Handler) handleOverview(
	ctx *gin.Context,
) ginx.Render {
	tenantID, err :=
		middleware.GetTenantIDFromContext(ctx)

	if err != nil {
		return ginx.Error(err)
	}

	stats, err := h.service.GetUserStats(
		ctx.Request.Context(),
		tenantID,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(
		srv.NewDashboardOverviewResponse(stats),
	)
}
