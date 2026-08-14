package param

import (
	//"encoding/json"

	"github.com/google/uuid"
)

type CreateBootstrapCredentialRequest struct {
	ResourceID uuid.UUID `json:"resource_id"`
}

type BindResourceIdentityRequest struct {
	ResourceID uuid.UUID `json:"resource_id"`

	IdentityType string `json:"identity_type"`

	HardwareID string `json:"hardware_id"`

	SerialNumber string `json:"serial_number"`
}

type ProvisionDeviceRequest struct {
	BootstrapToken string `json:"bootstrap_token"`

	HardwareID string `json:"hardware_id"`

	SerialNumber string `json:"serial_number"`

	CSR string `json:"csr"`
}

type AuthenticateDeviceRequest struct {
	CertificateFingerprint string `json:"certificate_fingerprint"`
}

type RenewCertificateRequest struct {
	CSR string `json:"csr"`
}

type RevokeCertificateRequest struct {
	ResourceID uuid.UUID `json:"resource_id"`

	CertificateID string `json:"certificate_id"`

	Reason string `json:"reason"`
}

type CertificateReq struct{
	Resource uuid.UUID `json:"resource_id" form:"resource_id"`
	CertificateID string `json:"certificate_id" form:"certificate_id"`
}