package api

import (
	"errors"
	"net/http"

	"github.com/OpenIndustrial/cloud/internal/kernel/security"
	"github.com/OpenIndustrial/cloud/internal/param"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SecurityHandler wraps the security service to expose its functionality via HTTP.
type SecurityHandler struct {
	service security.Service
}

// NewSecurityHandler creates a new security handler.
func NewSecurityHandler(service security.Service) *SecurityHandler {
	return &SecurityHandler{
		service: service,
	}
}

// RegisterSecurityRoutes registers all security-related HTTP routes.
func (h *SecurityHandler) RegisterSecurityRoutes(router *gin.RouterGroup) {
	// Endpoint for unauthenticated devices to get their initial identity.
	router.POST("/provision", h.provisionDevice)

	// Endpoint for internal services (e.g., MQTT broker) to authenticate a device.
	router.POST("/auth/device", h.authenticateDevice)

	// Endpoints for managing bootstrap credentials.
	credGroup := router.Group("/credentials")
	{
		credGroup.POST("/bootstrap", h.createBootstrapCredential)
		credGroup.POST("/:id/revoke", h.revokeCredential)
	}

	// Endpoints for managing certificates.
	certGroup := router.Group("/certificates")
	{
		certGroup.GET("/:id", h.getCertificate)
		//certGroup.POST("/:id/activate", h.activateCertificate)
		certGroup.POST("/identity", h.bindResourceIdentity)
		certGroup.POST("/:id/revoke", h.revokeCertificate)
		certGroup.POST("/:id/rotate", h.rotateCertificate)
	}

	// Endpoint for listing certificates for a specific resource.
	resourceGroup := router.Group("/resources")
	{
		resourceGroup.GET("/:id/certificates", h.listCertificates)
	}
}

// --- Request/Response Structs ---

type provisionDeviceRequest struct {
	BootstrapToken string `json:"bootstrap_token" binding:"required"`
	HardwareID     string `json:"hardware_id" binding:"required"`
	SerialNumber   string `json:"serial_number"`
	CSR            string `json:"csr" binding:"required"`
}

type authenticateDeviceRequest struct {
	CertificateFingerprint string `json:"certificate_fingerprint" binding:"required"`
}

type createBootstrapCredentialRequest struct {
	ResourceID string `json:"resource_id" binding:"required"`
}

type revokeCertificateRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type rotateCertificateRequest struct {
	CSR string `json:"csr" binding:"required"`
}

// --- Handler Implementations ---

func (h *SecurityHandler) provisionDevice(c *gin.Context) {
	var req param.ProvisionDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	resp, err := h.service.ProvisionDevice(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}		
	c.JSON(http.StatusOK, resp)
}

func (h *SecurityHandler) authenticateDevice(c *gin.Context) {
	var req param.AuthenticateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	resp, err := h.service.AuthenticateDevice(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SecurityHandler) createBootstrapCredential(c *gin.Context) {
	var req param.CreateBootstrapCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	resp, err := h.service.CreateBootstrapCredential(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *SecurityHandler) revokeCredential(c *gin.Context) {
	credID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credential id format"})
		return
	}

	if err := h.service.RevokeBootstrapCredential(c.Request.Context(), credID); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *SecurityHandler) getCertificate(c *gin.Context) {
	var req param.CertificateReq

	if err := c.ShouldBindQuery(req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return 
	}

	resp, err := h.service.GetCertificate(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SecurityHandler) listCertificates(c *gin.Context) {
	resID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id format"})
		return
	}

	resp, err := h.service.ListCertificates(c.Request.Context(), resID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// func (h *SecurityHandler) activateCertificate(c *gin.Context) {
// 	certID, err := uuid.Parse(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid certificate id format"})
// 		return
// 	}

// 	if err := h.service.ActivateCertificate(c.Request.Context(), certID); err != nil {
// 		h.handleError(c, err)
// 		return
// 	}
// 	c.Status(http.StatusNoContent)
// }

func (h *SecurityHandler) revokeCertificate(c *gin.Context) {
	// certID, err := uuid.Parse(c.Param("id"))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid certificate id format"})
	// 	return
	// }

	var req param.RevokeCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.service.RevokeCertificate(c.Request.Context(), req); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *SecurityHandler) rotateCertificate(c *gin.Context) {
	// certID, err := uuid.Parse(c.Param("id"))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid certificate id format"})
	// 	return
	// }

	var req param.RenewCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	resp, err := h.service.RenewCertificate(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SecurityHandler) bindResourceIdentity(c *gin.Context) {
	// resID, err := uuid.Parse(c.Param("resource_id"))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource_id format"})
	// 	return
	// }
	var req param.BindResourceIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, err)
		return
	}
	resp, err := h.service.BindResourceIdentity(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}


// handleError is a centralized error handler for the security service.
func (h *SecurityHandler) handleError(c *gin.Context, err error) {
	switch {
	// case errors.Is(err, security.ErrCredentialInvalid), errors.Is(err, security.ErrCertificateInvalid):
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, security.ErrCredentialConsumed):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, security.ErrCredentialRevoked), errors.Is(err, security.ErrCertificateRevoked):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, security.ErrIdentityMismatch):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, security.ErrCredentialNotFound), errors.Is(err, security.ErrIdentityNotFound), errors.Is(err, security.ErrCertificateNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		// For any other unexpected error, return a generic 500 error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an internal error occurred"})
	}
}

