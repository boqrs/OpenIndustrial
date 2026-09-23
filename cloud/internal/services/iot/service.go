package iot

import (
	"context"
	"errors"
	"strings"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
)

var (
	ErrInvalidAuthenticationRequest = errors.New(
		"invalid mqtt authentication request",
	)

	ErrDeviceNotFound = errors.New(
		"device not found",
	)
)

type service struct {
	repo     Repository
	security security.Service
}

// NewService creates the IoT service.
func NewService(
	repo Repository,
	securitySvc security.Service,
) Service {
	return &service{
		repo:     repo,
		security: securitySvc,
	}
}

// AuthenticateMQTT authenticates an MQTT client certificate.
//
// Authentication is delegated to Security:
//
//	Certificate fingerprint
//	        ↓
//	ResourceCertificate
//	        ↓
//	ResourceID
//
// IoT does not maintain another certificate registry.
func (s *service) AuthenticateMQTT(
	ctx context.Context,
	req AuthenticateMQTTRequest,
) (*MQTTAuthenticationResponse, error) {
	fingerprint := strings.TrimSpace(
		req.CertificateFingerprint,
	)

	if fingerprint == "" {
		return nil, ErrInvalidAuthenticationRequest
	}

	if s.security == nil {
		return nil, errors.New(
			"security service is not configured",
		)
	}

	auth, err := s.security.AuthenticateDevice(
		ctx,
		security.AuthenticateDeviceRequest{
			CertificateFingerprint: fingerprint,
		},
	)
	if err != nil {
		return nil, err
	}

	if auth == nil ||
		!auth.Authenticated ||
		auth.ResourceID == 0 {

		return nil, security.ErrCertificateNotFound
	}

	// Authentication proves that the certificate belongs to the
	// Resource. We still verify that the Resource has an actual
	// manufactured Device record.
	if s.repo == nil {
		return nil, errors.New(
			"iot repository is not configured",
		)
	}

	if _, err := s.repo.GetDeviceByResourceID(
		ctx,
		auth.ResourceID,
	); err != nil {
		return nil, err
	}

	return &MQTTAuthenticationResponse{
		Authenticated: true,
		ResourceID:    auth.ResourceID,
		CertificateID: auth.CertificateID,
	}, nil
}

// AuthorizeMQTT authorizes an already authenticated Resource.
func (s *service) AuthorizeMQTT(
	ctx context.Context,
	req AuthorizeMQTTRequest,
) error {
	_ = ctx

	if req.ResourceID == 0 {
		return ErrInvalidResourceID
	}

	if s.repo == nil {
		return errors.New(
			"iot repository is not configured",
		)
	}

	// Do not allow an arbitrary ResourceID supplied by an external
	// caller to become an MQTT identity. The Resource must exist as
	// a manufactured Device.
	if _, err := s.repo.GetDeviceByResourceID(
		ctx,
		req.ResourceID,
	); err != nil {
		return err
	}

	return AuthorizeMQTT(req)
}

// GetDeviceTopics returns the topic namespace of a Device.
func (s *service) GetDeviceTopics(
	ctx context.Context,
	resourceID uint,
) (*DeviceTopicsResponse, error) {
	if resourceID == 0 {
		return nil, ErrInvalidResourceID
	}

	if s.repo == nil {
		return nil, errors.New(
			"iot repository is not configured",
		)
	}

	if _, err := s.repo.GetDeviceByResourceID(
		ctx,
		resourceID,
	); err != nil {
		return nil, err
	}

	return GetDeviceTopics(
		resourceID,
	)
}

// DeviceOnline records an actual MQTT connection.
func (s *service) DeviceOnline(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {
	if resourceID == 0 {
		return nil, ErrInvalidResourceID
	}

	if s.repo == nil {
		return nil, errors.New(
			"iot repository is not configured",
		)
	}

	runtime, err := s.repo.SetOnline(
		ctx,
		resourceID,
	)
	if err != nil {
		return nil, err
	}

	return runtime, nil
}

// DeviceOffline records an MQTT disconnect.
func (s *service) DeviceOffline(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {
	if resourceID == 0 {
		return nil, ErrInvalidResourceID
	}

	if s.repo == nil {
		return nil, errors.New(
			"iot repository is not configured",
		)
	}

	runtime, err := s.repo.SetOffline(
		ctx,
		resourceID,
	)
	if err != nil {
		return nil, err
	}

	return runtime, nil
}

// DeviceHeartbeat refreshes the last-online timestamp.
func (s *service) DeviceHeartbeat(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {
	if resourceID == 0 {
		return nil, ErrInvalidResourceID
	}

	if s.repo == nil {
		return nil, errors.New(
			"iot repository is not configured",
		)
	}

	runtime, err := s.repo.UpdateLastOnline(
		ctx,
		resourceID,
	)
	if err != nil {
		return nil, err
	}

	return runtime, nil
}
