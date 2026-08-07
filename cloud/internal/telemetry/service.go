package telemetry

import (
	"context"
	"encoding/json"
	"time"
)

// MQTTMessage represents the payload from the gateway.
type MQTTMessage struct {
	Timestamp int64              `json:"time"`
	Metrics   map[string]float64 `json:"metrics"`
}

// Service provides business logic for handling telemetry data.
type Service struct {
	repo         Repository
	metricCache  map[string]string // A local cache for metric code -> metric ID
	ruleEngine   RuleEngine
}

// NewService creates a new Telemetry Service.
func NewService(repo Repository) *Service {
	return &Service{
		repo:        repo,
		metricCache: make(map[string]string),
		// ruleEngine would be initialized here
	}
}

// HandleMQTTMessage is the entry point for data coming from the MQTT consumer.
func (s *Service) HandleMQTTMessage(ctx context.Context, deviceID string, msg *MQTTMessage) error {
	// 1. Convert MQTT message to DataPoint objects.
	var points []*DataPoint
	timestamp := time.Unix(0, msg.Timestamp)
	for code, value := range msg.Metrics {
		points = append(points, &DataPoint{
			ResourceID: deviceID, // Assuming deviceID is the resourceID
			Metric:     code,
			Value:      value,
			Timestamp:  timestamp,
		})
	}

	if len(points) == 0 {
		return nil
	}

	// 2. Save the time-series data.
	if err := s.repo.SaveDataPoints(ctx, points); err != nil {
		return err
	}

	// 3. Update the real-time device state.
	stateValues, _ := json.Marshal(msg.Metrics)
	state := &DeviceState{
		DeviceID: deviceID,
		Status:   "online",
		LastSeen: timestamp,
		Values:   stateValues,
	}
	if err := s.repo.UpdateDeviceState(ctx, state); err != nil {
		// Log the error, but don't fail the whole operation.
	}

	// 4. Pass the data to the rule engine to check for alarms.
	// s.ruleEngine.Process(ctx, points)

	return nil
}

// RuleEngine is a placeholder for the alarm rule engine component.
type RuleEngine interface {
	Process(ctx context.Context, points []*DataPoint)
}