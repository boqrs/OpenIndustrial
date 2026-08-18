package security

import (
	"context"

	"github.com/OpenIndustrial/cloud/internal/kernel/security/provider"
)

// CertificateAuthorityAdapter implements the security.CertificateAuthority interface
// by wrapping and adapting a provider.CertificateAuthority instance.
type CertificateAuthorityAdapter struct {
	provider provider.CertificateAuthority
}

// NewCertificateAuthorityAdapter creates a new adapter that makes the provider's CA
// compatible with the interface required by the security service.
func NewCertificateAuthorityAdapter(p provider.CertificateAuthority) CertificateAuthority {
	return &CertificateAuthorityAdapter{provider: p}
}

// ValidateCSR adapts the call to the underlying provider, mapping only the fields
// that exist in the security.ParsedCSR struct.
func (a *CertificateAuthorityAdapter) ValidateCSR(csrPEM string) (*ParsedCSR, error) {
	providerCSR, err := a.provider.ValidateCSR(csrPEM)
	if err != nil {
		return nil, err
	}
	if providerCSR == nil {
		return nil, nil
	}

	// Explicitly map only the fields that exist in security.ParsedCSR.
	return &ParsedCSR{
		Subject:  providerCSR.Subject,
		URIs:     providerCSR.URIs,
		DNSNames: providerCSR.DNSNames,
	}, nil
}

// IssueCertificate adapts the call, converting request and response types.
// It does not pass ValidityDays, as this is handled by the provider's configuration.
func (a *CertificateAuthorityAdapter) IssueCertificate(ctx context.Context, req IssueCertificateRequest) (*IssuedCertificate, error) {
	// Convert security.IssueCertificateRequest to provider.IssueCertificateRequest.
	// Note that req.ValidityDays does not exist in the security layer request,
	// so we don't set it. The provider will use its default configured value.
	providerReq := provider.IssueCertificateRequest{
		ResourceID: req.ResourceID,
		CSR:        req.CSR,
	}

	// Call the wrapped provider
	providerCert, err := a.provider.IssueCertificate(ctx, providerReq)
	if err != nil {
		return nil, err
	}
	if providerCert == nil {
		return nil, nil
	}

	// Convert provider.IssuedCertificate back to security.IssuedCertificate
	return &IssuedCertificate{
		CertificateID:  providerCert.CertificateID,
		CertificatePEM: providerCert.CertificatePEM,
		Fingerprint:    providerCert.Fingerprint,
		SerialNumber:   providerCert.SerialNumber,
		Subject:        providerCert.Subject,
		Issuer:         providerCert.Issuer,
		NotBefore:      providerCert.NotBefore,
		NotAfter:       providerCert.NotAfter,
	}, nil
}

// RevokeCertificate implements the correct interface signature required by security.CertificateAuthority.
// It takes certificateID and reason as strings and converts them into a provider.RevokeCertificateRequest.
func (a *CertificateAuthorityAdapter) RevokeCertificate(ctx context.Context, certificateID string, reason string) error {
	// Convert the string reason to the provider's typed reason.
	// We will perform a simple mapping here.
	var providerReason provider.CertificateRevokeReason
	switch reason {
	case "KEY_COMPROMISE":
		providerReason = provider.CertificateRevokeReasonKeyCompromise
	case "CA_COMPROMISE":
		providerReason = provider.CertificateRevokeReasonCACompromise
	case "SUPERSEDED":
		providerReason = provider.CertificateRevokeReasonSuperseded
	case "CESSATION_OF_OPERATION":
		providerReason = provider.CertificateRevokeReasonCessationOfOperation
	case "PRIVILEGE_WITHDRAWN":
		providerReason = provider.CertificateRevokeReasonPrivilegeWithdrawn
	default:
		providerReason = provider.CertificateRevokeReasonUnspecified
	}

	// Construct the request object required by the provider.
	providerReq := provider.RevokeCertificateRequest{
		CertificateID: certificateID,
		// The provider's request also has a SerialNumber field, but the security
		// interface doesn't provide it here. We will leave it empty.
		// The underlying provider implementation (e.g., AWS) might only need the CertificateID.
		SerialNumber: "",
		Reason:       providerReason,
	}

	// Call the wrapped provider with the correct request object.
	return a.provider.RevokeCertificate(ctx, providerReq)
}