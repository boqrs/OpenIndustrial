package auth

import "time"

// DeviceCredential holds the authentication credentials for a device.
// This can be based on certificates or secure tokens.
type DeviceCredential struct {
	ID            string
	DeviceID      string // Links to the Device resource
	CertificateSN string
	SecretHash    string // For token-based authentication
	Status        string // e.g., "active", "revoked", "expired"
	CreatedAt     time.Time
}