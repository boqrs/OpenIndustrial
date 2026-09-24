package protocol

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidResourceID = errors.New(
	"invalid resource id",
)

// DeviceTopicPrefix returns:
//
// devices/{resourceID}
func DeviceTopicPrefix(
	resourceID uint,
) string {
	if resourceID == 0 {
		return ""
	}

	return fmt.Sprintf(
		"devices/%d",
		resourceID,
	)
}

// DeviceStatusTopic returns:
//
// devices/{resourceID}/status
//
// Device publishes.
// Cloud subscribes.
func DeviceStatusTopic(
	resourceID uint,
) string {
	if resourceID == 0 {
		return ""
	}

	return fmt.Sprintf(
		"%s/status",
		DeviceTopicPrefix(resourceID),
	)
}

// DeviceCommandTopic returns:
//
// devices/{resourceID}/command
//
// Cloud publishes.
// Device subscribes.
func DeviceCommandTopic(
	resourceID uint,
) string {
	if resourceID == 0 {
		return ""
	}

	return fmt.Sprintf(
		"%s/command",
		DeviceTopicPrefix(resourceID),
	)
}

// DeviceTopics contains the two MQTT business topics of a Device.
type DeviceTopics struct {
	ResourceID uint

	StatusTopic string

	CommandTopic string
}

// GetDeviceTopics returns all MQTT business topics of a Device.
func GetDeviceTopics(
	resourceID uint,
) (*DeviceTopics, error) {
	if resourceID == 0 {
		return nil, ErrInvalidResourceID
	}

	return &DeviceTopics{
		ResourceID: resourceID,

		StatusTopic: DeviceStatusTopic(
			resourceID,
		),

		CommandTopic: DeviceCommandTopic(
			resourceID,
		),
	}, nil
}

func IsDeviceStatusTopic(
	resourceID uint,
	topic string,
) bool {
	if resourceID == 0 {
		return false
	}

	return strings.TrimSpace(topic) ==
		DeviceStatusTopic(resourceID)
}

func IsDeviceCommandTopic(
	resourceID uint,
	topic string,
) bool {
	if resourceID == 0 {
		return false
	}

	return strings.TrimSpace(topic) ==
		DeviceCommandTopic(resourceID)
}
