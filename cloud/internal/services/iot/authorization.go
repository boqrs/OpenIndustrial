package iot

import (
	"errors"
	"strings"
)

var (
	ErrMQTTAccessDenied = errors.New(
		"mqtt access denied",
	)
)

// AuthorizeMQTT checks whether a Resource is allowed to perform
// an MQTT operation on the specified topic.
//
// Device permissions:
//
//	Publish:
//	  devices/{resource_id}/telemetry
//	  devices/{resource_id}/status
//	  devices/{resource_id}/event
//
//	Subscribe:
//	  devices/{resource_id}/command
//
// A Resource can never access another Resource's namespace.
func AuthorizeMQTT(
	req AuthorizeMQTTRequest,
) error {
	if req.ResourceID == 0 {
		return ErrInvalidResourceID
	}

	topic := strings.TrimSpace(
		req.Topic,
	)

	if topic == "" {
		return ErrMQTTAccessDenied
	}

	switch req.Action {

	case MQTTActionPublish:

		if isDeviceTelemetryTopic(
			req.ResourceID,
			topic,
		) {
			return nil
		}

		if isDeviceStatusTopic(
			req.ResourceID,
			topic,
		) {
			return nil
		}

		if isDeviceEventTopic(
			req.ResourceID,
			topic,
		) {
			return nil
		}

	case MQTTActionSubscribe:

		if isDeviceCommandTopic(
			req.ResourceID,
			topic,
		) {
			return nil
		}
	}

	return ErrMQTTAccessDenied
}
