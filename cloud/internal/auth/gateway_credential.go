package auth

import "time"

// GatewayCredential holds the authentication credentials for a gateway,
// typically based on X.509 certificates.
type GatewayCredential struct {
	ID              string
	GatewayID       string // Links to the Gateway resource
	CertificateSN   string // The serial number of the certificate
	Status          string // e.g., "active", "revoked", "expired"
	CreatedAt       time.Time
	ExpiresAt       time.Time
}