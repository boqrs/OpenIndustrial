package opcua

import "log"

// Publisher defines the interface for publishing samples to an external system.
type Publisher interface {
	Publish(sample Sample) error
}

// NoOpPublisher is a default publisher that does nothing.
// It's useful for testing or when no external publishing is configured.
type NoOpPublisher struct{}

// NewNoOpPublisher creates a new no-op publisher.
func NewNoOpPublisher() *NoOpPublisher {
	return &NoOpPublisher{}
}

// Publish logs the sample instead of sending it anywhere.
func (p *NoOpPublisher) Publish(sample Sample) error {
	log.Printf("NoOpPublisher: Publishing sample for ID %s, Value: %v", sample.ID, sample.Value)
	return nil
}