package gateway

import (
	"fmt"
	"sync"
)

// PublisherFactory is a function that creates a new Publisher instance.
type PublisherFactory func() Publisher

var (
	publisherFactories = make(map[string]PublisherFactory)
	publisherMu        sync.RWMutex
)

// RegisterPublisher registers a new publisher factory.
// It is intended to be called from the init function of each publisher implementation.
func RegisterPublisher(t string, factory PublisherFactory) {
	publisherMu.Lock()
	defer publisherMu.Unlock()
	if _, ok := publisherFactories[t]; ok {
		// Or log a warning, depending on desired behavior for duplicate registration
		panic(fmt.Sprintf("publisher type '%s' is already registered", t))
	}
	publisherFactories[t] = factory
}

// NewPublisher creates a new publisher instance of a given type.
func NewPublisher(t string) (Publisher, error) {
	publisherMu.RLock()
	defer publisherMu.RUnlock()
	factory, ok := publisherFactories[t]
	if !ok {
		return nil, fmt.Errorf("unknown publisher type '%s'", t)
	}
	return factory(), nil
}