package iot

import (
	"fmt"
	"strings"
)

// Topic represents a parsed MQTT topic for gateway events.
type Topic struct {
	GatewayID string
}

// ParseGatewayEventTopic parses a string into a Topic object.
// Format: oi/gateway/{gatewayId}/event
func ParseGatewayEventTopic(topic string) (Topic, error) {
	parts := strings.Split(topic, "/")
	if len(parts) != 4 {
		return Topic{}, fmt.Errorf("invalid topic format: %s, expected 4 parts", topic)
	}
	if parts[0] != "oi" || parts[1] != "gateway" || parts[3] != "event" {
		return Topic{}, fmt.Errorf("invalid topic structure: %s", topic)
	}

	return Topic{
		GatewayID: parts[2],
	}, nil
}