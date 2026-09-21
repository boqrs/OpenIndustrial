package security

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource"
)

var (
	ErrResourceNotFound      = errors.New("resource not found")
	ErrCredentialNotFound    = errors.New("credential not found")
	ErrCredentialInvalid     = errors.New("invalid credential")
	ErrCredentialRevoked     = errors.New("credential revoked")
	ErrCredentialConsumed    = errors.New("credential already consumed")
	ErrIdentityNotFound      = errors.New("resource identity not found")
	ErrIdentityMismatch      = errors.New("resource identity mismatch")
	ErrIdentityAlreadyExists = errors.New("resource identity already exists")
	ErrCertificateNotFound   = errors.New("certificate not found")
	ErrCertificateMismatch   = errors.New("certificate does not belong to resource")
	ErrCertificateRevoked    = errors.New("certificate revoked")
	ErrCertificateExpired    = errors.New("certificate expired")
)

type service struct {
	resources    resource.ResourceRepository
	credentials  CredentialRepository
	identities   IdentityRepository
	certificates CertificateRepository
	ca           CertificateAuthority
	mqtt         MQTTProvider
	uow          UnitOfWork
}

func NewService(
	resources resource.ResourceRepository,
	credentials CredentialRepository,
	identities IdentityRepository,
	certificates CertificateRepository,
	ca CertificateAuthority,
	mqtt MQTTProvider,
	uow UnitOfWork,
) Service {
	return &service{
		resources:    resources,
		credentials:  credentials,
		identities:   identities,
		certificates: certificates,
		ca:           ca,
		mqtt:         mqtt,
		uow:          uow,
	}
}

func (s *service) CreateBootstrapCredential(
	ctx context.Context,
	req CreateBootstrapCredentialRequest,
) (*BootstrapCredentialResponse, error) {
	if req.ResourceID == 0 {
		return nil, errors.New("resource_id is required")
	}

	exists, err := s.resources.Exists(
		ctx,
		req.ResourceID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"check resource: %w",
			err,
		)
	}

	if !exists {
		return nil, ErrResourceNotFound
	}

	secret, err := generateSecret(32)
	if err != nil {
		return nil, fmt.Errorf(
			"generate secret: %w",
			err,
		)
	}

	now := time.Now().UTC()

	credential := &model.ResourceCredential{
		ResourceID: req.ResourceID,
		Type:       model.CredentialTypeBootstrap,
		Status:     model.CredentialStatusActive,
		SecretHash: hashSecret(secret),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.credentials.Create(
		ctx,
		credential,
	); err != nil {
		return nil, fmt.Errorf(
			"create credential: %w",
			err,
		)
	}

	// The token format is:
	//
	//     <credential_id>.<secret>
	//
	// The database stores only SHA256(secret).
	//
	// The plaintext token is returned exactly once.
	token := formatBootstrapToken(
		credential.ID,
		secret,
	)

	return &BootstrapCredentialResponse{
		ResourceID:   req.ResourceID,
		CredentialID: credential.ID,
		Token:        token,
		CreatedAt:    now,
	}, nil
}

func (s *service) RevokeBootstrapCredential(
	ctx context.Context,
	resourceID uint,
) error {
	if resourceID == 0 {
		return errors.New("resource_id is required")
	}

	exists, err := s.resources.Exists(
		ctx,
		resourceID,
	)
	if err != nil {
		return fmt.Errorf(
			"check resource: %w",
			err,
		)
	}

	if !exists {
		return ErrResourceNotFound
	}

	if err := s.credentials.Revoke(
		ctx,
		resourceID,
	); err != nil {
		return fmt.Errorf(
			"revoke bootstrap credential: %w",
			err,
		)
	}

	return nil
}

func (s *service) BindResourceIdentity(
	ctx context.Context,
	req BindResourceIdentityRequest,
) (*ResourceIdentityResponse, error) {
	return s.bindResourceIdentity(ctx, req)
}

// BindResourceIdentityTx binds a canonical identity inside an already
// existing transaction.
//
// IMPORTANT:
// This method deliberately does not call UnitOfWork.Execute().
// The caller owns the transaction.
func (s *service) BindResourceIdentityTx(
	ctx context.Context,
	req BindResourceIdentityRequest,
) (*ResourceIdentityResponse, error) {
	return s.bindResourceIdentity(ctx, req)
}

func (s *service) bindResourceIdentity(
	ctx context.Context,
	req BindResourceIdentityRequest,
) (*ResourceIdentityResponse, error) {
	if req.ResourceID == 0 {
		return nil, errors.New(
			"resource_id is required",
		)
	}

	// A canonical identity must have at least one physical identity value.
	//
	// SerialNumber is the normal manufacturing identity.
	// HardwareID is optional because not every product necessarily exposes
	// one during manufacturing.
	if req.HardwareID == "" &&
		req.SerialNumber == "" {
		return nil, errors.New(
			"hardware_id or serial_number is required",
		)
	}

	exists, err := s.resources.Exists(
		ctx,
		req.ResourceID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"check resource: %w",
			err,
		)
	}

	if !exists {
		return nil, ErrResourceNotFound
	}

	// A resource may have only one canonical identity record.
	existing, err := s.identities.GetByResourceID(
		ctx,
		req.ResourceID,
	)

	if err != nil &&
		!errors.Is(err, ErrIdentityNotFound) {
		return nil, fmt.Errorf(
			"get identity: %w",
			err,
		)
	}

	if existing != nil {
		return nil, ErrIdentityAlreadyExists
	}

	// HardwareID must not belong to another resource.
	if req.HardwareID != "" {
		exists, err := s.identities.HardwareIDExists(
			ctx,
			req.HardwareID,
			&req.ResourceID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"check hardware id uniqueness: %w",
				err,
			)
		}

		if exists {
			return nil, ErrIdentityAlreadyExists
		}
	}

	tenantID := tenantIDFromContext(ctx)

	// SerialNumber must not belong to another resource.
	if req.SerialNumber != "" {
		exists, err := s.identities.SerialNumberExists(
			ctx,
			tenantID,
			req.SerialNumber,
			&req.ResourceID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"check serial number uniqueness: %w",
				err,
			)
		}

		if exists {
			return nil, ErrIdentityAlreadyExists
		}
	}

	now := time.Now().UTC()

	identity := &model.ResourceIdentity{
		ResourceID:   req.ResourceID,
		HardwareID:   req.HardwareID,
		SerialNumber: req.SerialNumber,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.identities.Create(
		ctx,
		identity,
	); err != nil {
		return nil, fmt.Errorf(
			"create resource identity: %w",
			err,
		)
	}

	return &ResourceIdentityResponse{
		ResourceID:   identity.ResourceID,
		IdentityType: req.IdentityType,
		HardwareID:   identity.HardwareID,
		SerialNumber: identity.SerialNumber,
		CreatedAt:    identity.CreatedAt,
	}, nil
}

func (s *service) ProvisionDevice(
	ctx context.Context,
	req ProvisionDeviceRequest,
) (*ProvisionDeviceResponse, error) {
	if req.BootstrapToken == "" {
		return nil, errors.New(
			"bootstrap_token is required",
		)
	}

	if req.HardwareID == "" &&
		req.SerialNumber == "" {
		return nil, errors.New(
			"hardware_id or serial_number is required",
		)
	}

	if req.CSR == "" {
		return nil, errors.New(
			"csr is required",
		)
	}

	credentialID, secret, err := parseBootstrapToken(
		req.BootstrapToken,
	)
	if err != nil {
		return nil, ErrCredentialInvalid
	}

	// If the request still supplies ID, validate it against the token.
	//
	// The token itself is authoritative, so callers do not need to
	// separately supply an ID.
	if req.ID != 0 && req.ID != credentialID {
		return nil, ErrCredentialInvalid
	}

	var result *ProvisionDeviceResponse

	err = s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			credential, err := s.credentials.GetForUpdate(
				txCtx,
				credentialID,
			)
			if err != nil {
				if errors.Is(
					err,
					ErrCredentialNotFound,
				) {
					return ErrCredentialNotFound
				}

				return fmt.Errorf(
					"get credential: %w",
					err,
				)
			}

			if credential.Type !=
				model.CredentialTypeBootstrap {
				return ErrCredentialInvalid
			}

			if credential.Status ==
				model.CredentialStatusRevoked {
				return ErrCredentialRevoked
			}

			if credential.Status ==
				model.CredentialStatusConsumed {
				return ErrCredentialConsumed
			}

			if !verifySecret(
				secret,
				credential.SecretHash,
			) {
				return ErrCredentialInvalid
			}

			resourceID := credential.ResourceID

			identity, err := s.identities.GetByResourceID(
				txCtx,
				resourceID,
			)
			if err != nil {
				if errors.Is(
					err,
					ErrIdentityNotFound,
				) {
					return ErrIdentityNotFound
				}

				return fmt.Errorf(
					"get resource identity: %w",
					err,
				)
			}

			if identity.HardwareID != "" &&
				identity.HardwareID != req.HardwareID {
				return ErrIdentityMismatch
			}

			if identity.SerialNumber != "" &&
				req.SerialNumber != "" &&
				identity.SerialNumber != req.SerialNumber {
				return ErrIdentityMismatch
			}

			csr, err := s.ca.ValidateCSR(req.CSR)
			if err != nil {
				return fmt.Errorf(
					"validate csr: %w",
					err,
				)
			}

			if err := validateCSRForResource(
				csr,
				resourceID,
			); err != nil {
				return err
			}

			issued, err := s.ca.IssueCertificate(
				txCtx,
				IssueCertificateRequest{
					ResourceID: resourceID,
					CSR:        req.CSR,
				},
			)
			if err != nil {
				return fmt.Errorf(
					"issue certificate: %w",
					err,
				)
			}

			now := time.Now().UTC()

			certificate := &model.ResourceCertificate{
				ResourceID:              resourceID,
				CertificateID:           issued.CertificateID,
				CertificateSerialNumber: issued.SerialNumber,
				Fingerprint:             issued.Fingerprint,
				Subject:                 issued.Subject,
				Issuer:                  issued.Issuer,
				Status:                  model.CertificateActive,
				NotBefore:               issued.NotBefore,
				NotAfter:                issued.NotAfter,
				CreatedAt:               now,
				UpdatedAt:               now,
			}

			if err := s.certificates.Create(
				txCtx,
				certificate,
			); err != nil {
				_ = s.ca.RevokeCertificate(
					ctx,
					issued.CertificateID,
					issued.SerialNumber,
					"database persistence failure",
				)

				return fmt.Errorf(
					"save certificate: %w",
					err,
				)
			}

			credential.Status =
				model.CredentialStatusConsumed
			credential.ConsumedAt = &now
			credential.UpdatedAt = now

			if err := s.credentials.Update(
				txCtx,
				credential,
			); err != nil {
				_ = s.ca.RevokeCertificate(
					ctx,
					issued.CertificateID,
					issued.SerialNumber,
					"credential consumption failure",
				)

				return fmt.Errorf(
					"consume credential: %w",
					err,
				)
			}

			result = &ProvisionDeviceResponse{
				ResourceID: resourceID,

				Certificate: CertificateResponse{
					ID:            certificate.ID,
					ResourceID:    resourceID,
					CertificateID: certificate.CertificateID,
					Fingerprint:   certificate.Fingerprint,
					Status:        string(certificate.Status),
					NotBefore:     certificate.NotBefore,
					NotAfter:      certificate.NotAfter,
					CreatedAt:     certificate.CreatedAt,
				},

				MQTT: MQTTConnectionInfo{
					Endpoint: s.mqtt.Endpoint(),
					Port:     s.mqtt.Port(),
					Protocol: s.mqtt.Protocol(),
					ClientID: strconv.Itoa(int(resourceID)),
				},

				ProvisionedAt: now,
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) AuthenticateDevice(
	ctx context.Context,
	req AuthenticateDeviceRequest,
) (*DeviceAuthenticationResponse, error) {
	if req.CertificateFingerprint == "" {
		return nil, errors.New(
			"certificate_fingerprint is required",
		)
	}

	certificate, err :=
		s.certificates.GetByFingerprint(
			ctx,
			req.CertificateFingerprint,
		)

	if err != nil {
		if errors.Is(
			err,
			ErrCertificateNotFound,
		) {
			return nil, ErrCertificateNotFound
		}

		return nil, fmt.Errorf(
			"get certificate: %w",
			err,
		)
	}

	if certificate.Status ==
		model.CertificateRevoked {
		return nil, ErrCertificateRevoked
	}

	now := time.Now().UTC()

	if now.After(certificate.NotAfter) {
		return nil, ErrCertificateExpired
	}

	if now.Before(certificate.NotBefore) {
		return nil, errors.New(
			"certificate is not yet valid",
		)
	}

	return &DeviceAuthenticationResponse{
		Authenticated: true,
		ResourceID:    certificate.ResourceID,
		CertificateID: certificate.CertificateID,
	}, nil
}

func (s *service) GetCertificate(
	ctx context.Context,
	req CertificateReq,
) (*model.ResourceCertificate, error) {
	if req.Resource == 0 {
		return nil, errors.New(
			"resource_id is required",
		)
	}

	if req.CertificateID == "" {
		return nil, errors.New(
			"certificate_id is required",
		)
	}

	certificate, err :=
		s.certificates.GetByCertificateID(
			ctx,
			req.CertificateID,
		)

	if err != nil {
		if errors.Is(
			err,
			ErrCertificateNotFound,
		) {
			return nil, ErrCertificateNotFound
		}

		return nil, fmt.Errorf(
			"get certificate: %w",
			err,
		)
	}

	if certificate.ResourceID != req.Resource {
		return nil, ErrCertificateMismatch
	}

	return certificate, nil
}

func (s *service) ListCertificates(
	ctx context.Context,
	resourceID uint,
) ([]model.ResourceCertificate, error) {
	if resourceID == 0 {
		return nil, errors.New(
			"resource_id is required",
		)
	}

	exists, err := s.resources.Exists(
		ctx,
		resourceID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"check resource: %w",
			err,
		)
	}

	if !exists {
		return nil, ErrResourceNotFound
	}

	certificates, err :=
		s.certificates.ListByResourceID(
			ctx,
			resourceID,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"list certificates: %w",
			err,
		)
	}

	return certificates, nil
}

func (s *service) RenewCertificate(
	ctx context.Context,
	req RenewCertificateRequest,
) (*model.ResourceCertificate, error) {
	if req.CSR == "" {
		return nil, errors.New(
			"csr is required",
		)
	}

	if req.ResourceID == 0 {
		return nil, errors.New(
			"resource_id is required",
		)
	}

	csr, err := s.ca.ValidateCSR(req.CSR)
	if err != nil {
		return nil, fmt.Errorf(
			"validate csr: %w",
			err,
		)
	}

	if err := validateCSRForResource(
		csr,
		req.ResourceID,
	); err != nil {
		return nil, err
	}

	issued, err := s.ca.IssueCertificate(
		ctx,
		IssueCertificateRequest{
			ResourceID: req.ResourceID,
			CSR:        req.CSR,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"issue certificate: %w",
			err,
		)
	}

	now := time.Now().UTC()

	certificate := &model.ResourceCertificate{
		ResourceID:              req.ResourceID,
		CertificateID:           issued.CertificateID,
		CertificateSerialNumber: issued.SerialNumber,
		Fingerprint:             issued.Fingerprint,
		Subject:                 issued.Subject,
		Issuer:                  issued.Issuer,
		Status:                  model.CertificateActive,
		NotBefore:               issued.NotBefore,
		NotAfter:                issued.NotAfter,
		CreatedAt:               now,
		UpdatedAt:               now,
	}

	if err := s.certificates.Create(
		ctx,
		certificate,
	); err != nil {
		_ = s.ca.RevokeCertificate(
			ctx,
			issued.CertificateID,
			issued.SerialNumber,
			"certificate persistence failed",
		)

		return nil, fmt.Errorf(
			"save renewed certificate: %w",
			err,
		)
	}

	return certificate, nil
}

func (s *service) RevokeCertificate(
	ctx context.Context,
	req RevokeCertificateRequest,
) error {
	if req.ResourceID == 0 {
		return errors.New(
			"resource_id is required",
		)
	}

	if req.CertificateID == "" {
		return errors.New(
			"certificate_id is required",
		)
	}

	certificate, err :=
		s.certificates.GetByCertificateID(
			ctx,
			req.CertificateID,
		)

	if err != nil {
		if errors.Is(
			err,
			ErrCertificateNotFound,
		) {
			return ErrCertificateNotFound
		}

		return fmt.Errorf(
			"get certificate: %w",
			err,
		)
	}

	if certificate.ResourceID !=
		req.ResourceID {
		return ErrCertificateMismatch
	}

	if certificate.Status ==
		model.CertificateRevoked {
		return nil
	}

	if err := s.ca.RevokeCertificate(
		ctx,
		certificate.CertificateID,
		certificate.CertificateSerialNumber,
		req.Reason,
	); err != nil {
		return fmt.Errorf(
			"revoke certificate from ca: %w",
			err,
		)
	}

	now := time.Now().UTC()

	certificate.Status = model.CertificateRevoked
	certificate.RevokedAt = &now
	certificate.UpdatedAt = now

	if err := s.certificates.Update(
		ctx,
		certificate,
	); err != nil {
		return fmt.Errorf(
			"update certificate: %w",
			err,
		)
	}

	return nil
}

// tenantIDFromContext returns the current tenant ID.
//
// ResourceIdentity currently has no tenant_id column, so the value is only
// used to satisfy the existing repository interface. The repository itself
// intentionally does not include tenant_id in its uniqueness query.
func tenantIDFromContext(ctx context.Context) (id [16]byte) {
	// Keep this helper local to avoid changing the existing security
	// repository contract in this patch.
	//
	// The repository currently ignores tenantID because ResourceIdentity
	// has no tenant_id field.
	//
	// We only need a zero UUID here.
	return id
}
