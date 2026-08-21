package provider

import (
	"time"

	"github.com/google/uuid"
)

// CertificateRevokeReason defines a provider-independent
// certificate revocation reason.
type CertificateRevokeReason string

const (
	CertificateRevokeReasonUnspecified CertificateRevokeReason =
		"unspecified"

	CertificateRevokeReasonKeyCompromise CertificateRevokeReason =
		"key_compromise"

	CertificateRevokeReasonCACompromise CertificateRevokeReason =
		"ca_compromise"

	CertificateRevokeReasonSuperseded CertificateRevokeReason =
		"superseded"

	CertificateRevokeReasonCessationOfOperation CertificateRevokeReason =
		"cessation_of_operation"

	CertificateRevokeReasonPrivilegeWithdrawn CertificateRevokeReason =
		"privilege_withdrawn"
	RevokeReasonKeyCompromise        CertificateRevokeReason = 
	    "KEY_COMPROMISE"

)

// ParsedCSR contains provider-independent information
// extracted from a PKCS#10 CSR.
type ParsedCSR struct {
	Subject string `json:"subject"`

	CommonName string `json:"common_name"`

	DNSNames []string `json:"dns_names"`

	IPAddresses []string `json:"ip_addresses"`

	URIs []string `json:"uris"`

	EmailAddresses []string `json:"email_addresses"`

	PublicKeyAlgorithm string `json:"public_key_algorithm"`

	PublicKeySize int `json:"public_key_size"`
}

// IssueCertificateRequest contains all information required
// by a CertificateAuthority to issue a certificate.
type IssueCertificateRequest struct {
	ResourceID uuid.UUID `json:"resource_id"`

	CSR string `json:"csr"`

	ValidityDays int `json:"validity_days"`
}

// IssuedCertificate contains provider-independent information
// about an issued X.509 certificate.
type IssuedCertificate struct {
	// CertificateID is an opaque provider-specific identifier.
	//
	// The upper layer MUST NOT parse this value.
	CertificateID string `json:"certificate_id"`

	// PEM encoded end-entity certificate.
	CertificatePEM string `json:"certificate_pem"`

	// SHA-256 fingerprint of the DER encoded certificate.
	Fingerprint string `json:"fingerprint"`

	// X.509 serial number.
	SerialNumber string `json:"serial_number"`

	Subject string `json:"subject"`

	Issuer string `json:"issuer"`

	NotBefore time.Time `json:"not_before"`

	NotAfter time.Time `json:"not_after"`
}

// RevokeCertificateRequest contains information required
// to revoke an issued certificate.
type RevokeCertificateRequest struct {
	// Opaque provider certificate identifier.
	CertificateID string `json:"certificate_id"`

	// X.509 certificate serial number.
	SerialNumber string `json:"serial_number"`

	Reason CertificateRevokeReason `json:"reason"`
}