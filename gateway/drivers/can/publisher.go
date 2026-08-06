package can

import "log"

// Publisher defines the interface for publishing signal values.
type Publisher interface {
	Publish(value SignalValue) error
}

// NoOpPublisher is a default publisher that does nothing but log the value.
// It's useful for testing or when no external publishing is configured.
type NoOpPublisher struct{}

// NewNoOpPublisher creates a new no-op publisher.
func NewNoOpPublisher() *NoOpPublisher {
	return &NoOpPublisher{}
}

// Publish logs the signal value instead of sending it anywhere.
func (p *NoOpPublisher) Publish(value SignalValue) error {
	log.Printf("NoOpPublisher: Publishing signal ID %s, Value: %v", value.ID, value.Value)
	return nil
}