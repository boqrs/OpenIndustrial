package protocol

import (
	"errors"
	"strings"
)

type MQTTAction string

const (
	MQTTActionPublish MQTTAction = "publish"

	MQTTActionSubscribe MQTTAction = "subscribe"
)

var ErrMQTTAccessDenied = errors.New(
	"mqtt access denied",
)

// AuthorizeMQTT defines the MQTT ACL from the Device perspective.
//
// Device:
//
//	publish   -> devices/{resourceID}/status
//	subscribe -> devices/{resourceID}/command
//
// Everything else is denied.
func AuthorizeMQTT(
	resourceID uint,
	action MQTTAction,
	topic string,
) error {
	if resourceID == 0 {
		return ErrInvalidResourceID
	}

	topic = strings.TrimSpace(topic)

	if topic == "" {
		return ErrMQTTAccessDenied
	}

	switch action {

	case MQTTActionPublish:

		if IsDeviceStatusTopic(
			resourceID,
			topic,
		) {
			return nil
		}

	case MQTTActionSubscribe:

		if IsDeviceCommandTopic(
			resourceID,
			topic,
		) {
			return nil
		}
	}

	return ErrMQTTAccessDenied
}
