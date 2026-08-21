package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// AliyunCAConfig contains Alibaba Cloud PCA configuration.
type AliyunCAConfig struct {
	RegionID string `json:"region_id"`

	CAInstanceID string `json:"ca_instance_id"`

	ValidityDays int `json:"validity_days"`
}

// AliyunCertificateClient is the minimum Alibaba Cloud
// PCA capability required by AliyunCA.
//
// The actual Alibaba Cloud SDK implementation is injected
// from outside this provider.
type AliyunCertificateClient interface {

	IssueCertificate(
		ctx context.Context,
		req AliyunIssueCertificateRequest,
	) (*AliyunIssueCertificateResponse, error)

	RevokeCertificate(
		ctx context.Context,
		req AliyunRevokeCertificateRequest,
	) error
}

type AliyunIssueCertificateRequest struct {
	CAInstanceID string

	CSR []byte

	ValidityDays int
}

type AliyunIssueCertificateResponse struct {
	CertificateID string

	CertificatePEM string
}

type AliyunRevokeCertificateRequest struct {
	CAInstanceID string

	CertificateID string

	SerialNumber string
}

// AliyunCA implements CertificateAuthority using
// Alibaba Cloud PCA.
type AliyunCA struct {
	client AliyunCertificateClient

	config AliyunCAConfig
}

// NewAliyunCA creates an Alibaba Cloud PCA provider.
func NewAliyunCA(
	client AliyunCertificateClient,
	config AliyunCAConfig,
) *AliyunCA {

	if config.ValidityDays <= 0 {
		config.ValidityDays = 365
	}

	return &AliyunCA{
		client: client,

		config: config,
	}
}

// ValidateCSR validates CSR locally.
func (p *AliyunCA) ValidateCSR(
	csrPEM string,
) (*ParsedCSR, error) {

	return ParseCSR(csrPEM)
}

// IssueCertificate issues an end-entity certificate
// using Alibaba Cloud PCA.
func (p *AliyunCA) IssueCertificate(
	ctx context.Context,
	req IssueCertificateRequest,
) (*IssuedCertificate, error) {

	if req.ResourceID == uuid.Nil {
		return nil, errors.New(
			"resource_id is required",
		)
	}

	if p.client == nil {
		return nil, errors.New(
			"aliyun certificate client is nil",
		)
	}

	if p.config.CAInstanceID == "" {
		return nil, errors.New(
			"aliyun ca instance id is required",
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

	result, err :=
		p.client.IssueCertificate(
			ctx,
			AliyunIssueCertificateRequest{

				CAInstanceID:
					p.config.CAInstanceID,

				CSR:
					csrDER,

				ValidityDays:
					validityDays,
			},
		)

	if err != nil {
		return nil, fmt.Errorf(
			"aliyun issue certificate: %w",
			err,
		)
	}

	if result == nil {
		return nil, errors.New(
			"aliyun returned nil certificate",
		)
	}

	if result.CertificateID == "" {
		return nil, errors.New(
			"aliyun returned empty certificate id",
		)
	}

	if result.CertificatePEM == "" {
		return nil, errors.New(
			"aliyun returned empty certificate pem",
		)
	}

	return ParseIssuedCertificate(
		result.CertificateID,
		result.CertificatePEM,
	)
}

// RevokeCertificate revokes a certificate
// from Alibaba Cloud PCA.
func (p *AliyunCA) RevokeCertificate(
	ctx context.Context,
	req RevokeCertificateRequest,
) error {

	if p.client == nil {
		return errors.New(
			"aliyun certificate client is nil",
		)
	}

	if p.config.CAInstanceID == "" {
		return errors.New(
			"aliyun ca instance id is required",
		)
	}

	if req.CertificateID == "" {
		return errors.New(
			"certificate_id is required",
		)
	}

	err :=
		p.client.RevokeCertificate(
			ctx,
			AliyunRevokeCertificateRequest{

				CAInstanceID:
					p.config.CAInstanceID,

				CertificateID:
					req.CertificateID,

				SerialNumber:
					req.SerialNumber,
			},
		)

	if err != nil {
		return fmt.Errorf(
			"aliyun revoke certificate: %w",
			err,
		)
	}

	return nil
}