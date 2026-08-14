package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OpenIndustrial/cloud/internal/kernel/resource"
	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"

	"github.com/google/uuid"
)

var (
	ErrResourceNotFound = errors.New(
		"resource not found",
	)

	ErrCredentialNotFound = errors.New(
		"credential not found",
	)

	ErrCredentialInvalid = errors.New(
		"invalid credential",
	)

	ErrCredentialRevoked = errors.New(
		"credential revoked",
	)

	ErrCredentialConsumed = errors.New(
		"credential already consumed",
	)

	ErrIdentityNotFound = errors.New(
		"resource identity not found",
	)

	ErrIdentityMismatch = errors.New(
		"resource identity mismatch",
	)

	ErrIdentityAlreadyExists = errors.New(
		"resource identity already exists",
	)

	ErrCertificateNotFound = errors.New(
		"certificate not found",
	)

	ErrCertificateMismatch = errors.New(
		"certificate does not belong to resource",
	)

	ErrCertificateRevoked = errors.New(
		"certificate revoked",
	)

	ErrCertificateExpired = errors.New(
		"certificate expired",
	)
)
/************************设备生命周期***********************************/
/*
                    工厂生产/烧录
                         │
                         ▼
              ┌────────────────────┐
              │   未激活设备       │
              │                    │
              │ Hardware ID        │
              │ Serial Number      │
              │ Bootstrap Token    │
              │ Device Private Key │
              └─────────┬──────────┘
                        │
                        │ ProvisionDevice
                        ▼
              ┌────────────────────┐
              │    已注册设备      │
              │                    │
              │ Resource ID        │
              │ Device Certificate│
              │ Private Key        │
              └─────────┬──────────┘
                        │
                        │ mTLS
                        ▼
                 ┌──────────────┐
                 │ AWS IoT Core │
                 └──────┬───────┘
                        │
                        │ MQTT
                        ▼
                   正常运行
                        │
              ┌─────────┴─────────┐
              │                   │
         Certificate            Certificate
           nearing                expired
           expiry                  │
              │                   │
              ▼                   ▼
     RenewCertificate          设备失效
              │
              ▼
          新证书
              │
              ▼
         继续 MQTT
*/
/***************************************************************/

type Service interface {
	// =========================================================
	// Bootstrap Credential
	// =========================================================

	CreateBootstrapCredential(ctx context.Context,req param.CreateBootstrapCredentialRequest) (*param.BootstrapCredentialResponse, error)
	RevokeBootstrapCredential(ctx context.Context,resourceID uuid.UUID) error

	// =========================================================
	// Resource Identity
	// =========================================================

	BindResourceIdentity(ctx context.Context,req param.BindResourceIdentityRequest) (*param.ResourceIdentityResponse, error)
	// =========================================================
	// Device Provisioning
	// =========================================================
	/**首次注册设备使用**/
	ProvisionDevice(ctx context.Context,req param.ProvisionDeviceRequest) (*param.ProvisionDeviceResponse, error)
	// =========================================================
	// Device Authentication
	// =========================================================

	AuthenticateDevice(ctx context.Context,req param.AuthenticateDeviceRequest) (*param.DeviceAuthenticationResponse, error)
	// =========================================================
	// Certificate
	// =========================================================

	GetCertificate(ctx context.Context, req param.CertificateReq) (*model.ResourceCertificate, error)
	ListCertificates(ctx context.Context,resourceID uuid.UUID) ([]model.ResourceCertificate, error)
	RenewCertificate(ctx context.Context,req param.RenewCertificateRequest) (*model.ResourceCertificate, error)
	RevokeCertificate(ctx context.Context,req param.RevokeCertificateRequest) error
}


type service struct {
	resources resource.ResourceRepository

	credentials CredentialRepository

	identities IdentityRepository

	certificates CertificateRepository

	ca CertificateAuthority

	mqtt MQTTProvider

	tx TransactionManager
}

func NewService(resources resource.ResourceRepository,credentials CredentialRepository,identities IdentityRepository,certificates CertificateRepository,ca CertificateAuthority,mqtt MQTTProvider,tx TransactionManager) Service {

	return &service{
		resources: resources,

		credentials: credentials,

		identities: identities,

		certificates: certificates,

		ca: ca,

		mqtt: mqtt,

		tx: tx,
	}
}

func (s *service) CreateBootstrapCredential(ctx context.Context,req param.CreateBootstrapCredentialRequest) (*param.BootstrapCredentialResponse, error) {
	if req.ResourceID == uuid.Nil {
		return nil, errors.New(
			"resource_id is required",
		)
	}

	exists, err := s.resources.Exists(ctx,req.ResourceID)
	if err != nil {
		return nil, fmt.Errorf(
			"check resource: %w",
			err,
		)
	}

	if !exists {
		return nil, ErrResourceNotFound
	}

	credentialID := uuid.New()

	secret, err := generateSecret(32)
	if err != nil {
		return nil, fmt.Errorf(
			"generate secret: %w",
			err,
		)
	}

	token := credentialID.String() + "." + secret

	now := time.Now().UTC()

	credential := &model.ResourceCredential{
		ID: credentialID,

		ResourceID: req.ResourceID,

		Type: model.CredentialTypeBootstrap,

		Status: model.CredentialStatusActive,

		SecretHash: hashSecret(secret),

		CreatedAt: now,

		UpdatedAt: now,
	}

	if err := s.credentials.Create(ctx,credential); err != nil {

		return nil, fmt.Errorf(
			"create credential: %w",
			err,
		)
	}

	return &param.BootstrapCredentialResponse{
		ResourceID: req.ResourceID,

		CredentialID: credentialID,

		Token: token,

		CreatedAt: now,
	}, nil
}

func (s *service) RevokeBootstrapCredential(ctx context.Context,resourceID uuid.UUID,
) error {

	if resourceID == uuid.Nil {
		return errors.New(
			"resource_id is required",
		)
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

	if err := s.credentials.Revoke(ctx, resourceID); err != nil {

		return fmt.Errorf(
			"revoke bootstrap credential: %w",
			err,
		)
	}

	return nil
}

func (s *service) BindResourceIdentity(ctx context.Context,req param.BindResourceIdentityRequest) (*param.ResourceIdentityResponse, error) {

	if req.ResourceID == uuid.Nil {
		return nil, errors.New(
			"resource_id is required",
		)
	}

	if req.IdentityType == "" {
		return nil, errors.New(
			"identity_type is required",
		)
	}

	if req.HardwareID == "" {
		return nil, errors.New(
			"hardware_id is required",
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

	existing, err := s.identities.GetByResourceID(
		ctx,
		req.ResourceID,
	)

	if err != nil {
		if !errors.Is(err, ErrIdentityNotFound) {
			return nil, fmt.Errorf(
				"get identity: %w",
				err,
			)
		}
	}

	if existing != nil {
		return nil, ErrIdentityAlreadyExists
	}

	now := time.Now().UTC()

	identity := &model.ResourceIdentity{
		ID: uuid.New(),

		ResourceID: req.ResourceID,

		//IdentityType: req.IdentityType,

		HardwareID: req.HardwareID,

		SerialNumber: req.SerialNumber,

		CreatedAt: now,

		UpdatedAt: now,
	}

	if err := s.identities.Create(ctx,identity); err != nil {
		return nil, fmt.Errorf(
			"create resource identity: %w",
			err,
		)
	}

	return &param.ResourceIdentityResponse{
		ResourceID: identity.ResourceID,

		//IdentityType: identity.IdentityType,

		HardwareID: identity.HardwareID,

		SerialNumber: identity.SerialNumber,

		CreatedAt: identity.CreatedAt,
	}, nil
}

func (s *service) ProvisionDevice(ctx context.Context,req param.ProvisionDeviceRequest) (*param.ProvisionDeviceResponse, error) {

	if req.BootstrapToken == "" {
		return nil, errors.New(
			"bootstrap_token is required",
		)
	}

	if req.HardwareID == "" {
		return nil, errors.New(
			"hardware_id is required",
		)
	}

	if req.CSR == "" {
		return nil, errors.New(
			"csr is required",
		)
	}

	credentialID, secret, err := parseBootstrapToken(req.BootstrapToken)
	if err != nil {
		return nil, ErrCredentialInvalid
	}

	var result *param.ProvisionDeviceResponse

	err = s.tx.WithinTransaction(ctx,func(txCtx context.Context) error {

			credential, err := s.credentials.GetForUpdate(txCtx,credentialID)
			if err != nil {
				if errors.Is(err,ErrCredentialNotFound) {
					return ErrCredentialNotFound
				}

				return fmt.Errorf(
					"get credential: %w",
					err,
				)
			}

			if credential.Type != model.CredentialTypeBootstrap {
				return ErrCredentialInvalid
			}

			if credential.Status == model.CredentialStatusRevoked {

				return ErrCredentialRevoked
			}

			if credential.Status == model.CredentialStatusConsumed {

				return ErrCredentialConsumed
			}

			if !verifySecret(secret,credential.SecretHash) {
				return ErrCredentialInvalid
			}

			resourceID := credential.ResourceID

			identity, err := s.identities.GetByResourceID(txCtx,resourceID)
			if err != nil {
				if errors.Is(err,ErrIdentityNotFound) {
					return ErrIdentityNotFound
				}

				return fmt.Errorf(
					"get resource identity: %w",
					err,
				)
			}

			if identity.HardwareID !=req.HardwareID {
				return ErrIdentityMismatch
			}

			if identity.SerialNumber != "" &&identity.SerialNumber != req.SerialNumber {

				return ErrIdentityMismatch
			}

			csr, err := s.ca.ValidateCSR(req.CSR)
			if err != nil {
				return fmt.Errorf(
					"validate csr: %w",
					err,
				)
			}

			if err := validateCSRForResource(csr,resourceID); err != nil {

				return err
			}

			issued, err := s.ca.IssueCertificate(txCtx,IssueCertificateRequest{ResourceID: resourceID,CSR: req.CSR})
			if err != nil {
				return fmt.Errorf(
					"issue certificate: %w",
					err,
				)
			}

			now := time.Now().UTC()

			certificate := &model.ResourceCertificate{
				ID: uuid.New(),

				ResourceID: resourceID,

				CertificateID: issued.CertificateID,

				Fingerprint: issued.Fingerprint,

				Subject: issued.Subject,

				Issuer: issued.Issuer,

				Status: model.CertificateActive,

				NotBefore: issued.NotBefore,

				NotAfter: issued.NotAfter,

				CreatedAt: now,

				UpdatedAt: now,
			}

			if err := s.certificates.Create(txCtx,certificate); err != nil {

				_ = s.ca.RevokeCertificate(
					ctx,
					issued.CertificateID,
					"database persistence failure",
				)

				return fmt.Errorf(
					"save certificate: %w",
					err,
				)
			}

			credential.Status =
				model.CredentialStatusConsumed

			credential.CreatedAt = time.Now()

			credential.UpdatedAt = now

			if err := s.credentials.Update(
				txCtx,
				credential,
			); err != nil {

				_ = s.ca.RevokeCertificate(
					ctx,
					issued.CertificateID,
					"credential consumption failure",
				)

				return fmt.Errorf(
					"consume credential: %w",
					err,
				)
			}

			result = &param.ProvisionDeviceResponse{
				ResourceID: resourceID,

				Certificate: param.CertificateResponse{
					ID: certificate.ID,

					ResourceID: resourceID,

					CertificateID:
						certificate.CertificateID,

					Fingerprint:
						certificate.Fingerprint,

					Status:
						string(certificate.Status),

					NotBefore:
						certificate.NotBefore,

					NotAfter:
						certificate.NotAfter,

					CreatedAt:
						certificate.CreatedAt,
				},

				MQTT: param.MQTTConnectionInfo{
					Endpoint: s.mqtt.Endpoint(),

					Port: s.mqtt.Port(),

					Protocol: s.mqtt.Protocol(),

					ClientID: resourceID.String(),
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

func (s *service) AuthenticateDevice(ctx context.Context,req param.AuthenticateDeviceRequest) (*param.DeviceAuthenticationResponse, error) {

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

	return &param.DeviceAuthenticationResponse{
		Authenticated: true,

		ResourceID: certificate.ResourceID,

		CertificateID:
			certificate.CertificateID,
	}, nil
}

func (s *service) GetCertificate(ctx context.Context, req param.CertificateReq) (*model.ResourceCertificate, error) {

	if req.Resource == uuid.Nil {
		return nil, errors.New(
			"resource_id is required",
		)
	}

	if req.CertificateID == "" {
		return nil, errors.New(
			"certificate_id is required",
		)
	}

	certificate, err := s.certificates.GetByCertificateID(ctx, req.CertificateID)
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

func (s *service) ListCertificates(ctx context.Context,resourceID uuid.UUID) ([]model.ResourceCertificate, error) {

	if resourceID == uuid.Nil {
		return nil, errors.New(
			"resource_id is required",
		)
	}

	exists, err :=
		s.resources.Exists(
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

	certificates, err := s.certificates.ListByResourceID(ctx,resourceID)

	if err != nil {
		return nil, fmt.Errorf(
			"list certificates: %w",
			err,
		)
	}

	return certificates, nil
}

func (s *service) RenewCertificate(ctx context.Context,req param.RenewCertificateRequest,) (*model.ResourceCertificate, error) {

	if req.CSR == "" {
		return nil, errors.New(
			"csr is required",
		)
	}

	resourceID, ok :=
		resourceIDFromContext(ctx)

	if !ok {
		return nil, errors.New(
			"authenticated resource identity missing",
		)
	}

	csr, err :=
		s.ca.ValidateCSR(
			req.CSR,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"validate csr: %w",
			err,
		)
	}

	if err :=
		validateCSRForResource(
			csr,
			resourceID,
		); err != nil {

		return nil, err
	}

	issued, err :=s.ca.IssueCertificate(ctx,IssueCertificateRequest{
				ResourceID: resourceID,
				CSR: req.CSR,
			},
		)

	if err != nil {
		return nil, fmt.Errorf(
			"issue certificate: %w",
			err,
		)
	}

	now := time.Now().UTC()

	certificate :=&model.ResourceCertificate{
			ID: uuid.New(),
			ResourceID:resourceID,
			CertificateID:issued.CertificateID,
			Fingerprint:issued.Fingerprint,
			Subject:issued.Subject,
			Issuer:issued.Issuer,
			Status:model.CertificateActive,
			NotBefore:issued.NotBefore,
			NotAfter:issued.NotAfter,
			CreatedAt:now,
			UpdatedAt:now,
		}

	if err :=s.certificates.Create(ctx,certificate); err != nil {

		_ = s.ca.RevokeCertificate(
			ctx,
			issued.CertificateID,
			"certificate persistence failed",
		)

		return nil, fmt.Errorf(
			"save renewed certificate: %w",
			err,
		)
	}

	return certificate, nil
}

func (s *service) RevokeCertificate(ctx context.Context,req param.RevokeCertificateRequest) error {

	if req.ResourceID == uuid.Nil {
		return errors.New(
			"resource_id is required",
		)
	}

	if req.CertificateID == "" {
		return errors.New(
			"certificate_id is required",
		)
	}

	certificate, err :=s.certificates.GetByCertificateID(ctx,req.CertificateID)

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

	if err :=s.ca.RevokeCertificate(
			ctx,
			certificate.CertificateID,
			req.Reason,
		); err != nil {

		return fmt.Errorf(
			"revoke certificate from ca: %w",
			err,
		)
	}

	now := time.Now().UTC()

	certificate.Status =model.CertificateRevoked
	certificate.RevokedAt =&now
	certificate.UpdatedAt =now

	if err :=s.certificates.Update(ctx,certificate); err != nil {

		return fmt.Errorf(
			"update certificate: %w",
			err,
		)
	}

	return nil
}

