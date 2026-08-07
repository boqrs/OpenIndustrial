package certificate

import "time"

// CertificateStatus defines the lifecycle status of a certificate.
type CertificateStatus string

const (
	StatusPending  CertificateStatus = "pending"
	StatusActive   CertificateStatus = "active"
	StatusRevoked  CertificateStatus = "revoked"
	StatusExpired  CertificateStatus = "expired"
)

// OwnerType indicates the type of entity that owns the certificate.
type OwnerType string

const (
	OwnerTypeGateway  OwnerType = "gateway"
	OwnerTypeProduct  OwnerType = "product"
)

// Certificate represents an X.509 certificate used for authenticating devices and gateways.
type Certificate struct {
	// ID is the unique identifier for the certificate, often the serial number.
	ID string

	// OwnerType specifies the kind of entity this certificate belongs to (e.g., "gateway", "product").
	OwnerType OwnerType

	// OwnerID is the ID of the entity that owns this certificate.
	OwnerID string

	// Status indicates the current state of the certificate in its lifecycle.
	Status CertificateStatus

	// PEM is the certificate data in PEM format.
	PEM string

	// ExpiresAt is the timestamp when the certificate will expire.
	ExpiresAt time.Time

	// CreatedAt is the timestamp when the certificate was issued.
	CreatedAt time.Time
}