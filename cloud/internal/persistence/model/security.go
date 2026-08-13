package model

import (
	"time"

	"github.com/google/uuid"
)

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

// ResourceCredential stores a credential used for authenticating a resource, typically for bootstrapping.
type ResourceCredential struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	// ResourceID is the foreign key linking this credential to its owner resource.
	ResourceID uuid.UUID `gorm:"type:uuid;not null;index"`
	Type       CredentialType   `gorm:"type:varchar(50);not null"`
	Status     CredentialStatus `gorm:"type:varchar(50);not null"`
	// SecretHash stores the hashed version of the secret (e.g., SHA256 of the bootstrap token).
	SecretHash string           `gorm:"type:varchar(255);not null"`
	CreatedAt  time.Time        `gorm:"autoCreateTime"`
	ConsumedAt *time.Time
	RevokedAt  *time.Time
	UpdatedAt  time.Time        `gorm:"autoUpdateTime"`
}

// ResourceIdentity stores the intrinsic, often hardware-based, identifiers of a resource.
type ResourceIdentity struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	// ResourceID is the foreign key linking this identity to its owner resource. It is unique
	// as a resource should only have one canonical identity record.
	ResourceID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	HardwareID   string    `gorm:"type:varchar(255);index"`
	SerialNumber string    `gorm:"type:varchar(255);index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// CertificateStatus defines the lifecycle status of a device certificate.
type CertificateStatus string

const (
	CertificatePending CertificateStatus = "pending"
	CertificateActive  CertificateStatus = "active"
	CertificateRevoked CertificateStatus = "revoked"
	CertificateExpired CertificateStatus = "expired"
)

// ResourceCertificate stores information about a X.509 certificate associated with a resource.
type ResourceCertificate struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	// ResourceID is the foreign key linking this certificate to its owner resource.
	ResourceID              uuid.UUID         `gorm:"type:uuid;not null;index"`
	CertificateID           string            `gorm:"type:varchar(255);not null;uniqueIndex"`
	CertificateSerialNumber string            `gorm:"type:varchar(255);index"`
	Fingerprint             string            `gorm:"type:varchar(255);not null;uniqueIndex"`
	Subject                 string            `gorm:"type:text"`
	Issuer                  string            `gorm:"type:text"`
	Status                  CertificateStatus `gorm:"type:varchar(50);not null"`
	NotBefore               time.Time
	NotAfter                time.Time
	CreatedAt               time.Time         `gorm:"autoCreateTime"`
	UpdatedAt               time.Time         `gorm:"autoUpdateTime"`
	ActivatedAt             *time.Time
	RevokedAt               *time.Time
}