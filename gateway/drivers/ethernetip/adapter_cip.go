package ethernetip

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// cipAdapter implements the Adapter interface using Explicit Messaging (CIP).
type cipAdapter struct {
	mu      sync.RWMutex
	session *Session
	client  *CIPClient
	config  ConnectionConfig
}

// NewCIPAdapter creates a new adapter for Explicit Messaging.
func NewCIPAdapter() Adapter {
	return &cipAdapter{}
}

func (a *cipAdapter) Connect(ctx context.Context, cfg ConnectionConfig) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.session != nil && a.IsConnected() {
		return ErrAlreadyConnected
	}

	a.config = cfg
	s := NewSession(cfg)
	if err := s.Connect(); err != nil {
		return fmt.Errorf("failed to establish session: %w", err)
	}

	a.session = s
	a.client = NewCIPClient(s)
	return nil
}

func (a *cipAdapter) Disconnect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.session == nil {
		return nil
	}
	err := a.session.Close()
	a.session = nil
	a.client = nil
	return err
}

func (a *cipAdapter) IsConnected() bool {
	// A simple check. A real implementation might involve a keep-alive or periodic check.
	return a.session != nil && a.session.conn != nil
}

func (a *cipAdapter) Read(ctx context.Context, points []PointMapping) ([]Sample, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.IsConnected() {
		return nil, ErrNotConnected
	}

	samples := make([]Sample, 0, len(points))
	for _, point := range points {
		var rawData []byte
		var dataType DataType
		var err error

		// Read data based on point type
		switch point.Type {
		case PointTypePLCTag:
			rawData, dataType, err = a.client.ReadTag(point.Tag)
			// If the point definition has a data type, it overrides the one from the PLC.
			if point.DataType != "" {
				dataType = point.DataType
			}
		case PointTypeCIPObject:
			// TODO: Implement Read CIP Object for Phase 2
			err = fmt.Errorf("reading CIP objects not implemented yet")
		default:
			err = fmt.Errorf("unknown point type: %s", point.Type)
		}

		// Process the result
		if err != nil {
			// Create a "bad quality" sample to indicate the read failure for this point.
			samples = append(samples, Sample{
				PointID:   point.ID,
				Timestamp: time.Now(),
				Quality:   QualityBad,
				Source:    "ethernetip",
			})
			// Log the error but continue to the next point.
			// log.Printf("Failed to read point %s (%s): %v", point.ID, point.Tag, err)
			continue
		}

		// Decode the raw data into a Go type.
		value, err := Decode(rawData, dataType)
		if err != nil {
			samples = append(samples, Sample{
				PointID:   point.ID,
				Timestamp: time.Now(),
				Quality:   QualityBad,
				Source:    "ethernetip",
			})
			// log.Printf("Failed to decode point %s (%s): %v", point.ID, point.Tag, err)
			continue
		}

		// Create a "good quality" sample.
		samples = append(samples, Sample{
			PointID:   point.ID,
			Value:     value,
			Timestamp: time.Now(),
			Quality:   QualityGood,
			Source:    "ethernetip",
		})
	}

	return samples, nil
}

func (a *cipAdapter) Write(ctx context.Context, point PointMapping, value interface{}) error {
	// TODO: Implement Write for Phase 1.
	// 1. Use Encode() to convert the value to []byte.
	// 2. Create a "Write Tag Service" CIP message.
	// 3. Send it using the client.
	return fmt.Errorf("write not implemented yet")
}

func (a *cipAdapter) Subscribe(ctx context.Context, mappings []PointMapping, ch chan<- Sample) error {
	// TODO: Implement Subscription for Phase 2.
	// For explicit mode, this would likely involve starting a dedicated polling goroutine
	// that calls Read() periodically and sends the results to the channel.
	return fmt.Errorf("subscribe not implemented yet")
}