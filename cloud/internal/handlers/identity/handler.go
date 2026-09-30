package identity

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"

	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/identity"

	"github.com/boqrs/zeus/ginx"
)

type Handler struct {
	service srv.Service
	auth    middleware.Service
}

func NewIdentityHandler(
	service srv.Service,
	auth middleware.Service,
) *Handler {

	return &Handler{
		service: service,
		auth:    auth,
	}

}

func (h *Handler) RouterRegister(router ginx.ZeroGinRouter) {

	group := router.Group("/api/v1/external")

	group.Handle(http.MethodPost, "/login", h.handleLogin) //todo: 前端需要修改
	group.Handle(http.MethodPost, "/logout", h.handleLogout)
	group.Handle(http.MethodPost, "/refresh", h.handleRefreshToken)

	// invitation
	authGroup := router.Group("/api/v1/external")
	authGroup.Use(h.auth.Authenticate())
	authGroup.Handle(http.MethodPost, "/identity/invitations", h.handleInviteUser)
	authGroup.Handle(http.MethodPost, "/identity/invitations/accept", h.handleAcceptInvitation)

	// users
	users := authGroup.Group("/users")
	users.Handle(http.MethodGet, "/lists", h.handleListUsers)
	users.Handle(http.MethodGet, "/:id", h.handleGetUser)
	users.Handle(http.MethodPut, "/:id", h.handleUpdateUser)
	users.Handle(http.MethodPost, "/:id/disable", h.handleDisableUser)
	users.Handle(http.MethodPost, "/:id/enable", h.handleEnableUser)
	users.Handle(http.MethodPost, "/:id/password", h.handleUpdatePassword)
	users.Handle(http.MethodPost, "/:id/reset-password", h.handleResetPassword)
}

// =====================================================
// Authentication
// =====================================================

func (h *Handler) handleLogin(ctx *gin.Context) ginx.Render {
	var req srv.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		return ginx.Error(
			fmt.Errorf("invalid request"),
		)

	}

	resp, err := h.service.Login(
		ctx.Request.Context(),
		req,
	)

	if err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(resp)

}

func (h *Handler) handleLogout(ctx *gin.Context) ginx.Render {

	var req srv.LogoutRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		return ginx.Error(err)

	}

	if err := h.service.Logout(
		ctx.Request.Context(),
		req,
	); err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(nil)

}

func (h *Handler) handleRefreshToken(ctx *gin.Context) ginx.Render {

	var req srv.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		return ginx.Error(err)

	}

	resp, err := h.service.RefreshToken(
		ctx.Request.Context(),
		req,
	)

	if err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(resp)

}

// =====================================================
// Invitation
// =====================================================

func (h *Handler) handleInviteUser(ctx *gin.Context) ginx.Render {

	var req srv.InviteUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		return ginx.Error(err)

	}

	if err := h.service.InviteUser(
		ctx.Request.Context(),
		req,
	); err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(nil)

}

func (h *Handler) handleAcceptInvitation(ctx *gin.Context) ginx.Render {

	var req srv.AcceptInvitationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		return ginx.Error(err)

	}

	if err := h.service.AcceptInvitation(
		ctx.Request.Context(),
		req,
	); err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(nil)

}

// =====================================================
// User
// =====================================================

func (h *Handler) handleListUsers(ctx *gin.Context) ginx.Render {

	TenantID, err := middleware.GetTenantIDFromContextV2(ctx)
	if err != nil {
		return ginx.Error(err)
	}
	req := srv.ListUsersRequest{
		TenantID: TenantID,
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {

		return ginx.Error(err)

	}

	users, err := h.service.ListUsers(
		ctx.Request.Context(),
		req,
	)

	if err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(users)

}

func (h *Handler) handleGetUser(ctx *gin.Context) ginx.Render {

	id, err := uuid.Parse(
		ctx.Param("id"),
	)

	if err != nil {

		return ginx.Error(err)

	}

	TenantID, err := middleware.GetTenantIDFromContextV2(ctx)
	if err != nil {
		return ginx.Error(err)
	}

	user, err := h.service.GetUser(ctx.Request.Context(), TenantID, id)

	if err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(user)

}

func (h *Handler) handleUpdateUser(ctx *gin.Context) ginx.Render {
	id, err := uuid.Parse(
		ctx.Param("id"),
	)

	if err != nil {

		return ginx.Error(err)

	}

	var req srv.UpdateUserRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {

		return ginx.Error(err)
	}

	req.TenantID, err = middleware.GetTenantIDFromContextV2(ctx)
	if err != nil {
		return ginx.Error(err)
	}

	req.UserID = id
	if err := h.service.UpdateUser(
		ctx.Request.Context(),
		req,
	); err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(nil)

}

func (h *Handler) handleDisableUser(ctx *gin.Context) ginx.Render {

	return h.changeUserStatus(
		ctx,
		true,
	)

}

func (h *Handler) handleEnableUser(ctx *gin.Context) ginx.Render {

	return h.changeUserStatus(
		ctx,
		false,
	)

}

func (h *Handler) changeUserStatus(ctx *gin.Context, disable bool) ginx.Render {

	id, err := uuid.Parse(
		ctx.Param("id"),
	)

	if err != nil {

		return ginx.Error(err)

	}

	tenantID, err := middleware.GetTenantIDFromContextV2(ctx)
	if err != nil {
		return ginx.Error(err)
	}

	if disable {
		err = h.service.DisableUser(
			ctx.Request.Context(),
			srv.DisableUserRequest{
				TenantID: tenantID,
				UserID:   id,
			},
		)

	} else {

		err = h.service.EnableUser(
			ctx.Request.Context(),
			srv.EnableUserRequest{
				TenantID: tenantID,
				UserID:   id,
			},
		)

	}

	if err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(nil)

}

// =====================================================
// Password
// =====================================================

func (h *Handler) handleUpdatePassword(ctx *gin.Context) ginx.Render {
	id, err := uuid.Parse(
		ctx.Param("id"),
	)

	if err != nil {

		return ginx.Error(err)

	}

	var req srv.UpdatePasswordRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {

		return ginx.Error(err)

	}

	req.TenantID, err = middleware.GetTenantIDFromContextV2(ctx)
	if err != nil {
		return ginx.Error(err)
	}

	req.UserID = id
	if err := h.service.UpdatePassword(
		ctx.Request.Context(),
		req,
	); err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(nil)

}

func (h *Handler) handleResetPassword(ctx *gin.Context) ginx.Render {

	id, err := uuid.Parse(
		ctx.Param("id"),
	)

	if err != nil {

		return ginx.Error(err)

	}

	var req srv.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		return ginx.Error(err)

	}

	req.TenantID, err = middleware.GetTenantIDFromContextV2(ctx)
	if err != nil {
		return ginx.Error(err)
	}

	req.UserID = id
	if err := h.service.ResetPassword(
		ctx.Request.Context(),
		req,
	); err != nil {

		return ginx.Error(err)

	}

	return ginx.Success(nil)

}
