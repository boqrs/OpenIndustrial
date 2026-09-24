package iot

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/google/uuid"
)

// CreateCommand creates a command record.
//
// Creating a command does not automatically send it.
func (s *service) CreateCommand(
	ctx context.Context,
	req *CreateCommandRequest,
) (*CommandResponse, error) {
	if req == nil ||
		req.DeviceID == 0 ||
		strings.TrimSpace(req.Command) == "" {

		return nil, ErrInvalidCommand
	}

	if s.repo == nil {
		return nil, errors.New(
			"iot repository is not configured",
		)
	}

	operatorID := pkg.UserIDFromContext(ctx)

	if operatorID == uuid.Nil {
		return nil, errors.New(
			"user id missing",
		)
	}

	device, err := s.repo.GetDeviceByID(
		ctx,
		req.DeviceID,
	)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, ErrDeviceNotFound
	}

	command := &model.DeviceCommand{
		DeviceID: req.DeviceID,

		OperatorID: operatorID,

		Command: strings.TrimSpace(
			req.Command,
		),

		Payload: req.Payload,

		Status: model.CommandStatusCreated,
	}

	if err := s.repo.CreateCommand(
		ctx,
		command,
	); err != nil {
		return nil, err
	}

	return commandToResponse(
		command,
	), nil
}

// SendCommand publishes an existing command to its Device.
func (s *service) SendCommand(
	ctx context.Context,
	commandID uint,
) (*CommandResponse, error) {
	if commandID == 0 {
		return nil, ErrCommandNotFound
	}

	if s.repo == nil {
		return nil, errors.New(
			"iot repository is not configured",
		)
	}

	if s.commandPublisher == nil {
		return nil, ErrMQTTCommandPublisherNotConfigured
	}

	command, err := s.repo.GetCommand(
		ctx,
		commandID,
	)
	if err != nil {
		return nil, err
	}

	if command == nil {
		return nil, ErrCommandNotFound
	}

	if command.Status != model.CommandStatusCreated {
		return nil, ErrInvalidCommandStatus
	}

	device, err := s.repo.GetDeviceByID(
		ctx,
		command.DeviceID,
	)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, ErrDeviceNotFound
	}

	messageID, err := s.commandPublisher.PublishCommand(
		ctx,
		device.ResourceID,
		command.ID,
		command.Command,
		command.Payload,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	command.Status =
		model.CommandStatusSent

	command.MessageID =
		messageID

	command.SentAt =
		&now

	if err := s.repo.UpdateCommand(
		ctx,
		command,
	); err != nil {
		return nil, err
	}

	return commandToResponse(
		command,
	), nil
}

// GetCommand returns a command by ID.
func (s *service) GetCommand(
	ctx context.Context,
	id uint,
) (*CommandResponse, error) {
	if id == 0 {
		return nil, ErrCommandNotFound
	}

	if s.repo == nil {
		return nil, errors.New(
			"iot repository is not configured",
		)
	}

	command, err := s.repo.GetCommand(
		ctx,
		id,
	)
	if err != nil {
		return nil, err
	}

	if command == nil {
		return nil, ErrCommandNotFound
	}

	return commandToResponse(
		command,
	), nil
}

// ListDeviceCommands returns all commands belonging to a Device.
func (s *service) ListDeviceCommands(
	ctx context.Context,
	deviceID uint,
) ([]*CommandResponse, error) {
	if deviceID == 0 {
		return nil, ErrInvalidCommand
	}

	if s.repo == nil {
		return nil, errors.New(
			"iot repository is not configured",
		)
	}

	commands, err := s.repo.ListDeviceCommands(
		ctx,
		deviceID,
	)
	if err != nil {
		return nil, err
	}

	result := make(
		[]*CommandResponse,
		0,
		len(commands),
	)

	for _, command := range commands {
		result = append(
			result,
			commandToResponse(command),
		)
	}

	return result, nil
}

// AcknowledgeCommand processes a command acknowledgement received
// from the Device.
func (s *service) AcknowledgeCommand(
	ctx context.Context,
	req *CommandAckRequest,
) error {
	if req == nil ||
		req.CommandID == 0 {

		return ErrCommandNotFound
	}

	if req.ResourceID == 0 {
		return ErrInvalidResourceID
	}

	if s.repo == nil {
		return errors.New(
			"iot repository is not configured",
		)
	}

	command, err := s.repo.GetCommand(
		ctx,
		req.CommandID,
	)
	if err != nil {
		return err
	}

	if command == nil {
		return ErrCommandNotFound
	}

	device, err := s.repo.GetDeviceByResourceID(
		ctx,
		req.ResourceID,
	)
	if err != nil {
		return err
	}

	if device == nil {
		return ErrDeviceNotFound
	}

	// A Device may only acknowledge its own commands.
	if device.ID != command.DeviceID {
		return ErrCommandDeviceMismatch
	}

	if command.Status != model.CommandStatusCreated &&
		command.Status != model.CommandStatusSent {

		return ErrInvalidCommandStatus
	}

	now := time.Now().UTC()

	if req.Success {
		command.Status =
			model.CommandStatusAcknowledged

		command.AcknowledgedAt =
			&now

		command.ErrorMessage = ""

	} else {
		command.Status =
			model.CommandStatusFailed

		command.ErrorMessage =
			strings.TrimSpace(req.Error)

		command.AcknowledgedAt =
			&now
	}

	return s.repo.UpdateCommand(
		ctx,
		command,
	)
}
