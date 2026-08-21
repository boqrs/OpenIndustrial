package provider

import "context"

// CertificateAuthority defines the provider-independent
// certificate authority capability.
//
// SecurityService depends only on this interface and does not
// know anything about AWS, Alibaba Cloud, PCA, ARN, etc.
type CertificateAuthority interface {
	ValidateCSR(csrPEM string) (*ParsedCSR, error)
	IssueCertificate(ctx context.Context,req IssueCertificateRequest) (*IssuedCertificate, error)
	RevokeCertificate(ctx context.Context,req RevokeCertificateRequest) error
}