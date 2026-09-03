package executors

import (
	"context"
	"errors"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
)

type CertificateIssueExecutor struct {
	ca security.CertificateAuthority
}

func NewCertificateIssueExecutor(
	ca security.CertificateAuthority,
) *CertificateIssueExecutor {
	return &CertificateIssueExecutor{
		ca: ca,
	}
}

func (e *CertificateIssueExecutor) Type() string {
	return OperationTypeCertificateIssue
}

func (e *CertificateIssueExecutor) Validate(
	ctx context.Context,
	input *OperationInput,
) error {
	if input == nil {
		return errors.New("operation input is nil")
	}

	if input.DeviceID == nil || *input.DeviceID == 0 {
		return errors.New("device ID is required")
	}

	if e.ca == nil {
		return errors.New("certificate authority is not configured")
	}

	csr, ok := input.Parameters["csr"].(string)
	if !ok || csr == "" {
		return errors.New("csr is required")
	}

	return nil
}

func (e *CertificateIssueExecutor) Execute(
	ctx context.Context,
	input *OperationInput,
) (*OperationOutput, error) {

	if err := e.Validate(ctx, input); err != nil {
		return nil, err
	}

	csr := input.Parameters["csr"].(string)

	issued, err := e.ca.IssueCertificate(
		ctx,
		security.IssueCertificateRequest{
			ResourceID: *input.DeviceID,
			CSR:        csr,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to issue device certificate: %w",
			err,
		)
	}

	return &OperationOutput{
		Result: map[string]any{
			"certificateId": issued.CertificateID,
			"certificate":   issued.CertificatePEM,
			"serialNumber":  issued.SerialNumber,
			"fingerprint":   issued.Fingerprint,
			"subject":       issued.Subject,
			"issuer":        issued.Issuer,
			"notBefore":     issued.NotBefore,
			"notAfter":      issued.NotAfter,
		},
		References: map[string]any{
			"certificateId": issued.CertificateID,
		},
	}, nil
}