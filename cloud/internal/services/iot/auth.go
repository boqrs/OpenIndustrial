package iot

import (
	"context"
	"errors"
	"strings"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/iot/protocol"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
)

// AuthenticateMQTT authenticates an MQTT client through Security.
//
// Certificate fingerprint
//
//	↓
//
// Security
//
//	↓
//
// ResourceID
//
//	↓
//
// Device
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

	if s.repo == nil {
		return nil, errors.New(
			"iot repository is not configured",
		)
	}

	// A certificate may identify a Resource, but MQTT IoT access
	// requires that the Resource actually represents a Device.
	if _, err := s.repo.GetDeviceByResourceID(
		ctx,
		auth.ResourceID,
	); err != nil {
		return nil, err
	}

	return &MQTTAuthenticationResponse{
		Authenticated: true,

		ResourceID: auth.ResourceID,

		CertificateID: auth.CertificateID,
	}, nil
}

// AuthorizeMQTT authorizes an authenticated Device to access
// its MQTT namespace.
func (s *service) AuthorizeMQTT(
	ctx context.Context,
	req AuthorizeMQTTRequest,
) error {
	if req.ResourceID == 0 {
		return ErrInvalidResourceID
	}

	if s.repo == nil {
		return errors.New(
			"iot repository is not configured",
		)
	}

	// ResourceID must resolve to an actual Device.
	if _, err := s.repo.GetDeviceByResourceID(
		ctx,
		req.ResourceID,
	); err != nil {
		return err
	}

	// Protocol layer owns the actual MQTT ACL rules.
	if err := protocol.AuthorizeMQTT(
		req.ResourceID,
		protocol.MQTTAction(req.Action),
		req.Topic,
	); err != nil {
		return ErrMQTTAccessDenied
	}

	return nil
}
