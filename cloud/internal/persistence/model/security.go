package model

import "time"

// CredentialType defines the type of a credential.
type CredentialType string

const (
	// CredentialTypeBootstrap is a one-time token used for initial device provisioning.
	CredentialTypeBootstrap CredentialType = "bootstrap"
)

// CredentialStatus defines the lifecycle status of a credential.
type CredentialStatus string

const (
	CredentialStatusActive   CredentialStatus = "active"
	CredentialStatusConsumed CredentialStatus = "consumed"
	CredentialStatusRevoked  CredentialStatus = "revoked"
)

// ResourceCredential stores a credential used for authenticating a resource,
// typically for bootstrapping.
type ResourceCredential struct {
	ID uint `gorm:"primaryKey"`

	// ResourceID links this credential to its owner resource.
	ResourceID uint `gorm:"not null;index"`

	Type   CredentialType   `gorm:"type:varchar(50);not null"`
	Status CredentialStatus `gorm:"type:varchar(50);not null"`

	// SecretHash stores the hashed version of the secret.
	SecretHash string `gorm:"type:varchar(255);not null"`

	CreatedAt  time.Time `gorm:"autoCreateTime"`
	ConsumedAt *time.Time
	RevokedAt  *time.Time
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (ResourceCredential) TableName() string {
	return "resource_credentials"
}

// ResourceIdentity stores the canonical physical identity of a resource.
//
// A resource has at most one canonical identity record.
// SerialNumber is normally the manufacturing identity.
// HardwareID is optional because not every product exposes one.
type ResourceIdentity struct {
	ID uint `gorm:"primaryKey"`

	ResourceID uint `gorm:"not null;index"`

	HardwareID   string `gorm:"type:varchar(255);index"`
	SerialNumber string `gorm:"type:varchar(255);index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (ResourceIdentity) TableName() string {
	return "resource_identities"
}

// CertificateStatus defines the lifecycle status of a device certificate.
type CertificateStatus string

const (
	CertificatePending CertificateStatus = "pending"
	CertificateActive  CertificateStatus = "active"
	CertificateRevoked CertificateStatus = "revoked"
	CertificateExpired CertificateStatus = "expired"
)

// ResourceCertificate stores information about an X.509 certificate
// associated with a resource.
type ResourceCertificate struct {
	ID uint `gorm:"primaryKey"`

	ResourceID uint `gorm:"not null;index"`

	// CertificateID is the opaque certificate identifier returned
	// by the external CA provider.
	//
	// Examples:
	//   AWS Private CA: certificate ARN
	//   Alibaba Cloud PCA: certificate ID
	//
	// The application must never assume this value is numeric.
	CertificateID string `gorm:"type:varchar(512);not null;index"`

	CertificateSerialNumber string `gorm:"type:varchar(255);index"`

	Fingerprint string `gorm:"type:varchar(255);not null;uniqueIndex"`

	Subject string `gorm:"type:text"`
	Issuer  string `gorm:"type:text"`

	Status CertificateStatus `gorm:"type:varchar(50);not null"`

	NotBefore time.Time
	NotAfter  time.Time

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	ActivatedAt *time.Time
	RevokedAt   *time.Time
}

func (ResourceCertificate) TableName() string {
	return "resource_certificates"
}
