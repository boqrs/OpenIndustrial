package api

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/boqrs/zeus/ginx"


	srv"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
)

// SecurityHandler wraps the security service to expose its functionality via HTTP.
type Handler struct {
	service srv.Service
}

// NewSecurityHandler creates a new security handler.
func NewHandler(service srv.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterSecurityRoutes registers all security-related HTTP routes.
func (h *Handler)  RouterRegister(router ginx.ZeroGinRouter) {
	// Endpoint for unauthenticated devices to get their initial identity.
	
	externalGroup := router.Group("/api/v1/external")
	externalGroup.Handle(http.MethodPost, "/provision", h.provisionDevice)
	externalGroup.Handle(http.MethodPost, "/auth/device", h.authenticateDevice)
	externalGroup.Handle(http.MethodPost, "/credentials/bootstrap", h.createBootstrapCredential)
	externalGroup.Handle(http.MethodPost, "/credentials/:id/revoke", h.revokeCredential)
	externalGroup.Handle(http.MethodGet, "/certificates/:id", h.getCertificate)
	externalGroup.Handle(http.MethodPost, "/certificates/identity", h.bindResourceIdentity)
	externalGroup.Handle(http.MethodPost, "/certificates/:id/revoke", h.revokeCertificate)
	externalGroup.Handle(http.MethodPost, "/certificates/:id/rotate", h.rotateCertificate)
	externalGroup.Handle(http.MethodPost, "/resources/:id/certificates", h.listCertificates)
}
// --- Handler Implementations ---

func (h *Handler) provisionDevice(ctx *gin.Context) ginx.Render{
	var req srv.ProvisionDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}
	resp, err := h.service.ProvisionDevice(ctx.Request.Context(), req)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(resp)		
}

func (h *Handler) authenticateDevice(ctx *gin.Context) ginx.Render{
	var req srv.AuthenticateDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}

	resp, err := h.service.AuthenticateDevice(ctx.Request.Context(), req)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(resp)
}

func (h *Handler) createBootstrapCredential(ctx *gin.Context) ginx.Render{
	var req srv.CreateBootstrapCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}

	resp, err := h.service.CreateBootstrapCredential(ctx.Request.Context(), req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) revokeCredential(ctx *gin.Context) ginx.Render{
	credID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid credential id format"))
	}

	if err := h.service.RevokeBootstrapCredential(ctx.Request.Context(), credID); err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(nil)
}

func (h *Handler) getCertificate(ctx *gin.Context) ginx.Render{
	var req srv.CertificateReq

	if err := ctx.ShouldBindQuery(req); err != nil{
		return ginx.Error(fmt.Errorf("invalid request body: %w", err.Error()))
	}

	resp, err := h.service.GetCertificate(ctx.Request.Context(), req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) listCertificates(ctx *gin.Context) ginx.Render{
	resID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ginx.Error(fmt.Errorf("invalid credential id format"))
	}

	resp, err := h.service.ListCertificates(ctx.Request.Context(), resID)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(resp)
}

func (h *Handler) revokeCertificate(ctx *gin.Context) ginx.Render{

	var req srv.RevokeCertificateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}

	if err := h.service.RevokeCertificate(ctx.Request.Context(), req); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

func (h *Handler) rotateCertificate(ctx *gin.Context) ginx.Render{
	var req srv.RenewCertificateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}

	resp, err := h.service.RenewCertificate(ctx.Request.Context(), req)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)

}

func (h *Handler) bindResourceIdentity(ctx *gin.Context) ginx.Render{
	var req srv.BindResourceIdentityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ginx.Error(fmt.Errorf("invalid param json"))
	}
	resp, err := h.service.BindResourceIdentity(ctx.Request.Context(), req)
	if err != nil {
		return ginx.Error(err)
	}
	return ginx.Success(resp)
}


