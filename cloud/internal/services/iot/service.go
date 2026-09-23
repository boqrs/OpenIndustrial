package iot

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
	"github.com/google/uuid"
)

var (
	ErrInvalidAuthenticationRequest = errors.New(
		"invalid mqtt authentication request",
	)

	ErrDeviceNotFound = errors.New(
		"device not found",
	)

	ErrCommandNotFound = errors.New(
		"device command not found",
	)

	ErrInvalidCommand = errors.New(
		"invalid device command",
	)

	ErrInvalidCommandStatus = errors.New(
		"invalid command status",
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

func (s *service) CreateCommand(
	ctx context.Context,
	req *CreateCommandRequest,
) (*CommandResponse, error) {

	if req == nil ||
		req.DeviceID == 0 ||
		req.Command == "" {

		return nil, ErrInvalidCommand

	}

	operatorID := pkg.UserIDFromContext(ctx)

	if operatorID == uuid.Nil {

		return nil,
			errors.New(
				"user id missing",
			)

	}

	entity := &model.DeviceCommand{

		DeviceID: req.DeviceID,

		OperatorID: operatorID,

		Command: req.Command,

		Payload: req.Payload,

		Status: model.CommandStatusCreated,
	}

	if err := s.repo.CreateCommand(
		ctx,
		entity,
	); err != nil {

		return nil, err

	}

	return commandToResponse(entity), nil

}

func (s *service) AcknowledgeCommand(
	ctx context.Context,
	req *CommandAckRequest,
) error {

	if req == nil ||
		req.CommandID == 0 {

		return ErrCommandNotFound

	}

	cmd, err := s.repo.GetCommand(
		ctx,
		req.CommandID,
	)

	if err != nil {

		return err

	}

	if cmd == nil {

		return ErrCommandNotFound

	}

	if cmd.Status != model.CommandStatusCreated &&
		cmd.Status != model.CommandStatusSent {

		return ErrInvalidCommandStatus

	}

	now := time.Now().UTC()

	if req.Success {

		cmd.Status =
			model.CommandStatusAcknowledged

		cmd.AcknowledgedAt = &now

	} else {

		cmd.Status =
			model.CommandStatusFailed

		cmd.ErrorMessage = req.Error

	}

	return s.repo.UpdateCommand(
		ctx,
		cmd,
	)

}

func (s *service) GetCommand(
	ctx context.Context,
	id uint,
) (
	*CommandResponse,
	error,
) {

	if id == 0 {

		return nil, ErrCommandNotFound

	}

	cmd, err := s.repo.GetCommand(
		ctx,
		id,
	)

	if err != nil {

		return nil, err

	}

	if cmd == nil {

		return nil, ErrCommandNotFound

	}

	return commandToResponse(cmd), nil

}

func (s *service) ListDeviceCommands(
	ctx context.Context,
	deviceID uint,
) (
	[]*CommandResponse,
	error,
) {

	if deviceID == 0 {

		return nil, ErrInvalidCommand

	}

	items, err := s.repo.ListDeviceCommands(
		ctx,
		deviceID,
	)

	if err != nil {

		return nil, err

	}

	result := make(
		[]*CommandResponse,
		0,
		len(items),
	)

	for _, item := range items {

		result = append(
			result,
			commandToResponse(item),
		)

	}

	return result, nil

}
