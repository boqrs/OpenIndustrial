package security


type CreateBootstrapCredentialRequest struct {
	ResourceID uint `json:"resource_id"`
}

type BindResourceIdentityRequest struct {
	ResourceID uint `json:"resource_id"`
	IdentityType string `json:"identity_type"`
	HardwareID string `json:"hardware_id"`
	SerialNumber string `json:"serial_number"`
}

//TODO: 需要确认这里是否合适
type ProvisionDeviceRequest struct {
	BootstrapToken string `json:"bootstrap_token"`
	HardwareID string `json:"hardware_id"`
	SerialNumber string `json:"serial_number"`
	CSR string `json:"csr"`
	ID uint  `json:"id"`
}

type AuthenticateDeviceRequest struct {
	CertificateFingerprint string `json:"certificate_fingerprint"`
}

type RenewCertificateRequest struct {
	CSR string `json:"csr"`
	ResourceID uint `json:"resource_id"`
}

type RevokeCertificateRequest struct {
	ResourceID uint `json:"resource_id"`
	CertificateID uint `json:"certificate_id"`
	Reason string `json:"reason"`
}

type CertificateReq struct{
	Resource uint `json:"resource_id" form:"resource_id"`
	CertificateID uint `json:"certificate_id" form:"certificate_id"`
}