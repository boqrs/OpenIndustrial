package security

import (
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)



type CertificateRepository interface {
    Create(ctx context.Context,certificate *model.ResourceCertificate) error
    GetActiveByResourceID(ctx context.Context,resourceID uuid.UUID) (*model.ResourceCertificate, error)
    GetByCertificateID(ctx context.Context,certificateID string) (*model.ResourceCertificate, error)
    ListByResourceID(ctx context.Context,resourceID uuid.UUID) ([]model.ResourceCertificate, error)
    Activate(ctx context.Context,id uuid.UUID,activatedAt time.Time) error
    Revoke(ctx context.Context,id uuid.UUID,revokedAt time.Time) error
	GetByFingerprint(ctx context.Context, fingerprint string) (*model.ResourceCertificate, error)
	Update(ctx context.Context, cert *model.ResourceCertificate) error

}

type IdentityRepository interface {
    GetByResourceID(ctx context.Context,resourceID uuid.UUID) (*model.ResourceIdentity, error)
    Create(ctx context.Context,identity *model.ResourceIdentity) error
    CreateOrUpdate(ctx context.Context,identity *model.ResourceIdentity) error
    HardwareIDExists(ctx context.Context,hardwareID string,excludeResourceID *uuid.UUID) (bool, error)
    SerialNumberExists(ctx context.Context,tenantID uuid.UUID,serialNumber string,excludeResourceID *uuid.UUID) (bool, error)
}

type CredentialRepository interface {
    Create(ctx context.Context,credential *model.ResourceCredential) error
    GetActive(ctx context.Context,resourceID uuid.UUID,credentialType model.CredentialType) (*model.ResourceCredential, error)
    GetByID(ctx context.Context,id uuid.UUID) (*model.ResourceCredential, error)
    Consume(ctx context.Context,id uuid.UUID,consumedAt time.Time) error
    Revoke(ctx context.Context,id uuid.UUID) error
	// GetForUpdate obtains a row with SELECT ... FOR UPDATE.
	//
	// This is important for bootstrap credential consumption,
	// because two devices must not be able to consume the same
	// credential concurrently.
	GetForUpdate(ctx context.Context,id uuid.UUID) (*model.ResourceCredential, error)
	Update(ctx context.Context, cred *model.ResourceCredential) error
}

type CertificateAuthority interface {
	ValidateCSR(csrPEM string) (*ParsedCSR, error)
	IssueCertificate(ctx context.Context,req IssueCertificateRequest) (*IssuedCertificate, error)
	RevokeCertificate(ctx context.Context,certificateID string,reason string) error
}

type ParsedCSR struct {
	Subject string

	URIs []string

	DNSNames []string
}

type IssueCertificateRequest struct {
	ResourceID uuid.UUID

	CSR string
}

type IssuedCertificate struct {
	CertificateID string

	CertificatePEM string

	SerialNumber string

	Fingerprint string

	Subject string

	Issuer string

	NotBefore time.Time

	NotAfter time.Time
}

type MQTTProvider interface {
	Endpoint() string
	Port() int
	Protocol() string
}

type TransactionManager interface {
	WithinTransaction(ctx context.Context,fn func(ctx context.Context) error) error
}

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

	CreateBootstrapCredential(ctx context.Context,req CreateBootstrapCredentialRequest) (*BootstrapCredentialResponse, error)
	RevokeBootstrapCredential(ctx context.Context,resourceID uuid.UUID) error

	// =========================================================
	// Resource Identity
	// =========================================================

	BindResourceIdentity(ctx context.Context,req BindResourceIdentityRequest) (*ResourceIdentityResponse, error)
	// =========================================================
	// Device Provisioning
	// =========================================================
	/**首次注册设备使用**/
	ProvisionDevice(ctx context.Context,req ProvisionDeviceRequest) (*ProvisionDeviceResponse, error)
	// =========================================================
	// Device Authentication
	// =========================================================

	AuthenticateDevice(ctx context.Context,req AuthenticateDeviceRequest) (*DeviceAuthenticationResponse, error)
	// =========================================================
	// Certificate
	// =========================================================

	GetCertificate(ctx context.Context, req CertificateReq) (*model.ResourceCertificate, error)
	ListCertificates(ctx context.Context,resourceID uuid.UUID) ([]model.ResourceCertificate, error)
	RenewCertificate(ctx context.Context,req RenewCertificateRequest) (*model.ResourceCertificate, error)
	RevokeCertificate(ctx context.Context,req RevokeCertificateRequest) error
}