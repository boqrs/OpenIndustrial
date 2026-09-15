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

func NewCAIssueExecutor(
	ca security.CertificateAuthority,
) *CertificateIssueExecutor {
	return &CertificateIssueExecutor{
		ca: ca,
	}
}

func (e *CertificateIssueExecutor) Type() string {
	return OperationTypeCAIssue
}

func (e *CertificateIssueExecutor) Validate(
	ctx context.Context,
	input *OperationInput,
) error {
	if input == nil {
		return errors.New("operation input is nil")
	}

	if e.ca == nil {
		return errors.New("certificate authority is not configured")
	}

	if input.ExecutionID == 0 {
		return errors.New("execution ID is required")
	}

	if input.ExecutionOperationID == 0 {
		return errors.New("execution operation ID is required")
	}

	if input.WorkOrderID == 0 {
		return errors.New("work order ID is required")
	}

	if input.Parameters == nil {
		return errors.New("operation parameters are required")
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
			ResourceID: input.WorkOrderID,
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
			"certificate_id": issued.CertificateID,
			"certificate":    issued.CertificatePEM,
			"serial_number":  issued.SerialNumber,
			"fingerprint":    issued.Fingerprint,
			"subject":        issued.Subject,
			"issuer":         issued.Issuer,
			"not_before":     issued.NotBefore,
			"not_after":      issued.NotAfter,
		},
		References: map[string]any{
			"certificate_id": issued.CertificateID,
		},
	}, nil
}
