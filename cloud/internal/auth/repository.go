package auth

import "context"

// Repository defines the interface for storing and retrieving authentication-related data.
type Repository interface {
	// Session management
	CreateSession(ctx context.Context, session *Session) error
	GetSession(ctx context.Context, id string) (*Session, error)
	DeleteSession(ctx context.Context, id string) error

	// Gateway credentials
	CreateGatewayCredential(ctx context.Context, cred *GatewayCredential) error
	GetGatewayCredentialBySN(ctx context.Context, sn string) (*GatewayCredential, error)

	// Device credentials
	CreateDeviceCredential(ctx context.Context, cred *DeviceCredential) error
	GetDeviceCredential(ctx context.Context, deviceID string) (*DeviceCredential, error)
}