package event

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// memoryBus is an in-memory implementation of the Bus interface.
type memoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]reflect.Value
}

// NewMemoryBus creates a new in-memory event bus.
func NewMemoryBus() Bus {
	return &memoryBus{
		handlers: make(map[string][]reflect.Value),
	}
}

// Publish sends an event to all registered handlers.
func (b *memoryBus) Publish(ctx context.Context, eventData interface{}) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	eventType := reflect.TypeOf(eventData).String()
	handlers, ok := b.handlers[eventType]
	if !ok {
		return nil
	}

	args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(eventData)}

	for _, handler := range handlers {
		go func(h reflect.Value) {
			ret := h.Call(args)
			if len(ret) > 0 && !ret[0].IsNil() {
				fmt.Printf("event handler for %s returned an error: %v\n", eventType, ret[0].Interface())
			}
		}(handler)
	}

	return nil
}

// Subscribe registers a handler for a specific event type.
func (b *memoryBus) Subscribe(eventType string, handler interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	hVal := reflect.ValueOf(handler)
	if hVal.Kind() != reflect.Func {
		return fmt.Errorf("handler must be a function")
	}

	b.handlers[eventType] = append(b.handlers[eventType], hVal)
	return nil
}