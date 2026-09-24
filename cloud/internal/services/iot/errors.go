package iot

import "errors"

var (
	ErrInvalidResourceID = errors.New(
		"invalid resource id",
	)

	ErrInvalidAuthenticationRequest = errors.New(
		"invalid mqtt authentication request",
	)

	ErrDeviceNotFound = errors.New(
		"device not found",
	)

	ErrInvalidCommand = errors.New(
		"invalid device command",
	)

	ErrCommandNotFound = errors.New(
		"device command not found",
	)

	ErrInvalidCommandStatus = errors.New(
		"invalid command status",
	)

	ErrCommandDeviceMismatch = errors.New(
		"command does not belong to device",
	)

	ErrMQTTAccessDenied = errors.New(
		"mqtt access denied",
	)

	ErrMQTTCommandPublisherNotConfigured = errors.New(
		"mqtt command publisher is not configured",
	)

	ErrMQTTStatusSubscriberNotConfigured = errors.New(
		"mqtt status subscriber is not configured",
	)
)
