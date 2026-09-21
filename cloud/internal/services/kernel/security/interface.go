package security

import (
	"context"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

type CertificateRepository interface {
	Create(
		ctx context.Context,
		certificate *model.ResourceCertificate,
	) error

	GetActiveByResourceID(
		ctx context.Context,
		resourceID uint,
	) (*model.ResourceCertificate, error)

	GetByCertificateID(
		ctx context.Context,
		certificateID string,
	) (*model.ResourceCertificate, error)

	ListByResourceID(
		ctx context.Context,
		resourceID uint,
	) ([]model.ResourceCertificate, error)

	Activate(
		ctx context.Context,
		id uint,
		activatedAt time.Time,
	) error

	Revoke(
		ctx context.Context,
		id uint,
		revokedAt time.Time,
	) error

	GetByFingerprint(
		ctx context.Context,
		fingerprint string,
	) (*model.ResourceCertificate, error)

	Update(
		ctx context.Context,
		cert *model.ResourceCertificate,
	) error
}

type IdentityRepository interface {
	GetByResourceID(
		ctx context.Context,
		resourceID uint,
	) (*model.ResourceIdentity, error)

	Create(
		ctx context.Context,
		identity *model.ResourceIdentity,
	) error

	CreateOrUpdate(
		ctx context.Context,
		identity *model.ResourceIdentity,
	) error

	HardwareIDExists(
		ctx context.Context,
		hardwareID string,
		excludeResourceID *uint,
	) (bool, error)

	SerialNumberExists(
		ctx context.Context,
		tenantID uuid.UUID,
		serialNumber string,
		excludeResourceID *uint,
	) (bool, error)
}

type CredentialRepository interface {
	Create(
		ctx context.Context,
		credential *model.ResourceCredential,
	) error

	GetActive(
		ctx context.Context,
		resourceID uint,
		credentialType model.CredentialType,
	) (*model.ResourceCredential, error)

	GetByID(
		ctx context.Context,
		id uint,
	) (*model.ResourceCredential, error)

	Consume(
		ctx context.Context,
		id uint,
		consumedAt time.Time,
	) error

	Revoke(
		ctx context.Context,
		id uint,
	) error

	GetForUpdate(
		ctx context.Context,
		id uint,
	) (*model.ResourceCredential, error)

	Update(
		ctx context.Context,
		cred *model.ResourceCredential,
	) error
}

type CertificateAuthority interface {
	ValidateCSR(
		csrPEM string,
	) (*ParsedCSR, error)

	IssueCertificate(
		ctx context.Context,
		req IssueCertificateRequest,
	) (*IssuedCertificate, error)

	RevokeCertificate(
		ctx context.Context,
		certificateID string,
		serialNumber string,
		reason string,
	) error
}

type ParsedCSR struct {
	Subject  string
	URIs     []string
	DNSNames []string
}

type IssueCertificateRequest struct {
	ResourceID uint
	CSR        string
}

type IssuedCertificate struct {
	CertificateID  string
	CertificatePEM string
	SerialNumber   string
	Fingerprint    string
	Subject        string
	Issuer         string
	NotBefore      time.Time
	NotAfter       time.Time
}

type MQTTProvider interface {
	Endpoint() string
	Port() int
	Protocol() string
}

type UnitOfWork interface {
	Execute(
		ctx context.Context,
		fn func(ctx context.Context) error,
	) error
}

type Service interface {
	CreateBootstrapCredential(
		ctx context.Context,
		req CreateBootstrapCredentialRequest,
	) (*BootstrapCredentialResponse, error)

	RevokeBootstrapCredential(
		ctx context.Context,
		resourceID uint,
	) error

	BindResourceIdentity(
		ctx context.Context,
		req BindResourceIdentityRequest,
	) (*ResourceIdentityResponse, error)

	BindResourceIdentityTx(
		ctx context.Context,
		req BindResourceIdentityRequest,
	) (*ResourceIdentityResponse, error)

	ProvisionDevice(
		ctx context.Context,
		req ProvisionDeviceRequest,
	) (*ProvisionDeviceResponse, error)

	AuthenticateDevice(
		ctx context.Context,
		req AuthenticateDeviceRequest,
	) (*DeviceAuthenticationResponse, error)

	GetCertificate(
		ctx context.Context,
		req CertificateReq,
	) (*model.ResourceCertificate, error)

	ListCertificates(
		ctx context.Context,
		resourceID uint,
	) ([]model.ResourceCertificate, error)

	RenewCertificate(
		ctx context.Context,
		req RenewCertificateRequest,
	) (*model.ResourceCertificate, error)

	RevokeCertificate(
		ctx context.Context,
		req RevokeCertificateRequest,
	) error
}
