package provider

import "time"

type CertificateRevokeReason string

const (
	CertificateRevokeReasonUnspecified CertificateRevokeReason = "unspecified"

	CertificateRevokeReasonKeyCompromise CertificateRevokeReason = "key_compromise"

	CertificateRevokeReasonCACompromise CertificateRevokeReason = "ca_compromise"

	CertificateRevokeReasonSuperseded CertificateRevokeReason = "superseded"

	CertificateRevokeReasonCessationOfOperation CertificateRevokeReason = "cessation_of_operation"

	CertificateRevokeReasonPrivilegeWithdrawn CertificateRevokeReason = "privilege_withdrawn"
)

type ParsedCSR struct {
	Subject            string   `json:"subject"`
	CommonName         string   `json:"common_name"`
	DNSNames           []string `json:"dns_names"`
	IPAddresses        []string `json:"ip_addresses"`
	URIs               []string `json:"uris"`
	EmailAddresses     []string `json:"email_addresses"`
	PublicKeyAlgorithm string   `json:"public_key_algorithm"`
	PublicKeySize      int      `json:"public_key_size"`
}

type IssueCertificateRequest struct {
	ResourceID   uint   `json:"resource_id"`
	CSR          string `json:"csr"`
	ValidityDays int    `json:"validity_days"`
}

type IssuedCertificate struct {
	// Opaque provider-specific certificate identifier.
	//
	// AWS:
	//   certificate ARN
	//
	// Alibaba Cloud:
	//   provider certificate ID
	//
	// The upper layer must never parse this value.
	CertificateID string `json:"certificate_id"`

	CertificatePEM string `json:"certificate_pem"`

	Fingerprint string `json:"fingerprint"`

	SerialNumber string `json:"serial_number"`

	Subject string `json:"subject"`

	Issuer string `json:"issuer"`

	NotBefore time.Time `json:"not_before"`

	NotAfter time.Time `json:"not_after"`
}

type RevokeCertificateRequest struct {
	// Opaque provider-specific certificate identifier.
	CertificateID string `json:"certificate_id"`

	// X.509 certificate serial number.
	SerialNumber string `json:"serial_number"`

	Reason CertificateRevokeReason `json:"reason"`
}
