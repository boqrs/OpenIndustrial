package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"
	"github.com/aws/aws-sdk-go-v2/service/acmpca/types"
)

// AWSCAConfig contains AWS Private CA configuration.
type AWSCAConfig struct {
	CAARN string `json:"ca_arn"`

	CertificateTemplateARN string `json:"certificate_template_arn"`

	ValidityDays int `json:"validity_days"`

	SigningAlgorithm string `json:"signing_algorithm"`
}

// AWSCA implements CertificateAuthority using AWS Private CA.
type AWSCA struct {
	client *acmpca.Client

	config AWSCAConfig
}

// NewAWSCA creates an AWS Private CA provider.
func NewAWSCA(
	cfg awssdk.Config,
	config AWSCAConfig,
) *AWSCA {

	if config.ValidityDays <= 0 {
		config.ValidityDays = 365
	}

	if config.SigningAlgorithm == "" {
		config.SigningAlgorithm =
			"SHA256WITHRSA"
	}

	return &AWSCA{
		client:
			acmpca.NewFromConfig(cfg),

		config:
			config,
	}
}

// ValidateCSR validates a CSR locally.
//
// No AWS API call is required.
func (p *AWSCA) ValidateCSR(
	csrPEM string,
) (*ParsedCSR, error) {

	return ParseCSR(csrPEM)
}

// IssueCertificate issues an end-entity certificate
// through AWS Private CA.
func (p *AWSCA) IssueCertificate(
	ctx context.Context,
	req IssueCertificateRequest,
) (*IssuedCertificate, error) {

	if req.ResourceID == 0 {
		return nil, errors.New(
			"resource_id is required",
		)
	}

	if p.config.CAARN == "" {
		return nil, errors.New(
			"aws ca arn is required",
		)
	}

	csrDER, err :=
		ParseCSRDER(req.CSR)

	if err != nil {
		return nil, err
	}

	validityDays :=
		req.ValidityDays

	if validityDays <= 0 {
		validityDays =
			p.config.ValidityDays
	}

	input := &acmpca.IssueCertificateInput{
		CertificateAuthorityArn:
			awssdk.String(
				p.config.CAARN,
			),

		Csr:
			csrDER,

		SigningAlgorithm:
			parseAWSSigningAlgorithm(
				p.config.SigningAlgorithm,
			),

		Validity: &types.Validity{
			Type:
				types.ValidityPeriodTypeDays,

			Value:
				awssdk.Int64(
					int64(validityDays),
				),
		},
	}

	if p.config.CertificateTemplateARN != "" {
		input.TemplateArn =
			awssdk.String(
				p.config.CertificateTemplateARN,
			)
	}

	result, err :=
		p.client.IssueCertificate(
			ctx,
			input,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"aws issue certificate: %w",
			err,
		)
	}

	if result.CertificateArn == nil ||
		*result.CertificateArn == "" {

		return nil, errors.New(
			"aws returned empty certificate arn",
		)
	}

	certificateID :=
		*result.CertificateArn

	certificatePEM, err :=
		p.waitCertificate(
			ctx,
			certificateID,
		)

	if err != nil {
		return nil, err
	}

	certificateIDUint, _ := strconv.Atoi(certificateID)

	return ParseIssuedCertificate(
		uint(certificateIDUint),
		certificatePEM,
	)
}

// waitCertificate waits until AWS Private CA
// makes the issued certificate available.
func (p *AWSCA) waitCertificate(
	ctx context.Context,
	certificateID string,
) (string, error) {

	const (
		maxAttempts = 20
		interval    = 500 * time.Millisecond
	)

	for i := 0; i < maxAttempts; i++ {

		result, err :=
			p.client.GetCertificate(
				ctx,
				&acmpca.GetCertificateInput{
					CertificateAuthorityArn:
						awssdk.String(
							p.config.CAARN,
						),

					CertificateArn:
						awssdk.String(
							certificateID,
						),
				},
			)

		if err == nil {

			if result.Certificate == nil ||
				*result.Certificate == "" {

				return "",
					errors.New(
						"aws returned empty certificate",
					)
			}

			return *result.Certificate, nil
		}

		select {

		case <-ctx.Done():
			return "", ctx.Err()

		case <-time.After(interval):
		}
	}

	return "",
		errors.New(
			"timeout waiting aws certificate",
		)
}

// RevokeCertificate revokes a certificate from AWS Private CA.
func (p *AWSCA) RevokeCertificate(
	ctx context.Context,
	req RevokeCertificateRequest,
) error {

	if req.CertificateID == 0 {
		return errors.New(
			"certificate_id is required",
		)
	}

	if req.SerialNumber == "" {
		return errors.New(
			"serial_number is required",
		)
	}

	if p.config.CAARN == "" {
		return errors.New(
			"aws ca arn is required",
		)
	}

	_, err :=
		p.client.RevokeCertificate(
			ctx,
			&acmpca.RevokeCertificateInput{

				CertificateAuthorityArn:
					awssdk.String(
						p.config.CAARN,
					),

				// CertificateArn:
				// 	awssdk.String(
				// 		req.CertificateID,
				// 	),

				CertificateSerial:
					awssdk.String(
						req.SerialNumber,
					),

				RevocationReason:
					mapAWSRevokeReason(
						req.Reason,
					),
			},
		)

	if err != nil {
		return fmt.Errorf(
			"aws revoke certificate: %w",
			err,
		)
	}

	return nil
}

func mapAWSRevokeReason(
	reason CertificateRevokeReason,
) types.RevocationReason {

	switch reason {

	case CertificateRevokeReasonKeyCompromise:
		return types.RevocationReasonKeyCompromise

	case CertificateRevokeReasonCACompromise:
		return types.RevocationReasonKeyCompromise

	case CertificateRevokeReasonSuperseded:
		return types.RevocationReasonSuperseded

	case CertificateRevokeReasonCessationOfOperation:
		return types.RevocationReasonCessationOfOperation

	case CertificateRevokeReasonPrivilegeWithdrawn:
		return types.RevocationReasonPrivilegeWithdrawn

	default:
		return types.RevocationReasonUnspecified
	}
}

func parseAWSSigningAlgorithm(
	value string,
) types.SigningAlgorithm {

	switch value {

	case "SHA256WITHRSA":
		return types.SigningAlgorithmSha256withrsa

	case "SHA384WITHRSA":
		return types.SigningAlgorithmSha384withrsa

	case "SHA512WITHRSA":
		return types.SigningAlgorithmSha512withrsa

	case "SHA256WITHECDSA":
		return types.SigningAlgorithmSha256withecdsa

	case "SHA384WITHECDSA":
		return types.SigningAlgorithmSha384withecdsa

	case "SHA512WITHECDSA":
		return types.SigningAlgorithmSha512withecdsa

	default:
		return types.SigningAlgorithmSha256withrsa
	}
}