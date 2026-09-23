package iot

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidResourceID = errors.New(
		"invalid resource id",
	)
)

// DeviceTopicPrefix returns the root MQTT namespace of a device.
func DeviceTopicPrefix(resourceID uint) string {
	if resourceID == 0 {
		return ""
	}

	return fmt.Sprintf(
		"devices/%d",
		resourceID,
	)
}

// DeviceTelemetryTopic returns the telemetry topic.
func DeviceTelemetryTopic(resourceID uint) string {
	return deviceTopic(
		resourceID,
		"telemetry",
	)
}

// DeviceStatusTopic returns the runtime status topic.
func DeviceStatusTopic(resourceID uint) string {
	return deviceTopic(
		resourceID,
		"status",
	)
}

// DeviceEventTopic returns the event topic.
func DeviceEventTopic(resourceID uint) string {
	return deviceTopic(
		resourceID,
		"event",
	)
}

// DeviceCommandTopic returns the command topic.
func DeviceCommandTopic(resourceID uint) string {
	return deviceTopic(
		resourceID,
		"command",
	)
}

// GetDeviceTopics returns the complete MQTT topic namespace
// belonging to a Resource.
func GetDeviceTopics(
	resourceID uint,
) (*DeviceTopicsResponse, error) {
	if resourceID == 0 {
		return nil, ErrInvalidResourceID
	}

	return &DeviceTopicsResponse{
		ResourceID: resourceID,

		TelemetryTopic: DeviceTelemetryTopic(
			resourceID,
		),

		StatusTopic: DeviceStatusTopic(
			resourceID,
		),

		EventTopic: DeviceEventTopic(
			resourceID,
		),

		CommandTopic: DeviceCommandTopic(
			resourceID,
		),
	}, nil
}

func deviceTopic(
	resourceID uint,
	name string,
) string {
	if resourceID == 0 {
		return ""
	}

	return fmt.Sprintf(
		"%s/%s",
		DeviceTopicPrefix(resourceID),
		name,
	)
}

func isDeviceTelemetryTopic(
	resourceID uint,
	topic string,
) bool {
	return strings.TrimSpace(topic) ==
		DeviceTelemetryTopic(resourceID)
}

func isDeviceStatusTopic(
	resourceID uint,
	topic string,
) bool {
	return strings.TrimSpace(topic) ==
		DeviceStatusTopic(resourceID)
}

func isDeviceEventTopic(
	resourceID uint,
	topic string,
) bool {
	return strings.TrimSpace(topic) ==
		DeviceEventTopic(resourceID)
}

func isDeviceCommandTopic(
	resourceID uint,
	topic string,
) bool {
	return strings.TrimSpace(topic) ==
		DeviceCommandTopic(resourceID)
}
