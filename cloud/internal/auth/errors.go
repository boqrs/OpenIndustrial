package auth

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrSessionNotFound      = errors.New("session not found")
	ErrTokenInvalid         = errors.New("token is invalid")
	ErrCertificateRevoked   = errors.New("certificate has been revoked")
	ErrCredentialNotFound   = errors.New("credential not found")
)