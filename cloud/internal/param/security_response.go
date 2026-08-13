package param

import (
	"time"

	"github.com/google/uuid"
)


type CertificateResponse struct {
	ID uuid.UUID `json:"id"`

	ResourceID uuid.UUID `json:"resource_id"`

	CertificateID string `json:"certificate_id"`

	Fingerprint string `json:"fingerprint"`

	Status string `json:"status"`

	NotBefore time.Time `json:"not_before"`

	NotAfter time.Time `json:"not_after"`

	CreatedAt time.Time `json:"created_at"`

	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

type DeviceAuthenticationResponse struct {
	Authenticated bool `json:"authenticated"`

	ResourceID uuid.UUID `json:"resource_id"`

	CertificateID string `json:"certificate_id"`
}

type ProvisionDeviceResponse struct {
	ResourceID uuid.UUID `json:"resource_id"`

	Certificate CertificateResponse `json:"certificate"`

	MQTT MQTTConnectionInfo `json:"mqtt"`

	ProvisionedAt time.Time `json:"provisioned_at"`
}

type MQTTConnectionInfo struct {
	Endpoint string `json:"endpoint"`

	Port int `json:"port"`

	Protocol string `json:"protocol"`

	ClientID string `json:"client_id"`
}

type ResourceIdentityResponse struct {
	ResourceID uuid.UUID `json:"resource_id"`

	IdentityType string `json:"identity_type"`

	HardwareID string `json:"hardware_id"`

	SerialNumber string `json:"serial_number"`

	CreatedAt time.Time `json:"created_at"`
}

type BootstrapCredentialResponse struct {
	ResourceID   uuid.UUID `json:"resource_id"`
	CredentialID uuid.UUID `json:"credential_id"`

	// 明文 Token 只返回一次。
	// 数据库绝对不保存这个值。
	Token string `json:"token"`

	CreatedAt time.Time `json:"created_at"`
}

