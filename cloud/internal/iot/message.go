package iot

import "time"

const (
	MessageGatewayOnline  = "gateway.online"
	MessageGatewayOffline = "gateway.offline"
	MessagePointChanged   = "point.changed"
	MessageDeviceAlarm    = "device.alarm"
)

// Message represents the unified data structure for all IoT messages from a gateway.
type Message struct {
	Type      string      `json:"type"`
	GatewayID string      `json:"gatewayId"`
	DeviceID  string      `json:"deviceId,omitempty"`
	PointID   string      `json:"pointId,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	Quality   uint8       `json:"quality,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}