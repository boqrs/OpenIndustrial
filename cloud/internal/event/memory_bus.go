package event

import (
	"log"
	"sync"
)

type memoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewMemoryBus() Bus {
	return &memoryBus{
		handlers: make(map[string][]Handler),
	}
}

func (b *memoryBus) Subscribe(
	eventName string,
	handler Handler,
) error {
	if eventName == "" {
		return ErrEventNameEmpty
	}
	if handler == nil {
		return ErrHandlerNil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventName] = append(
		b.handlers[eventName],
		handler,
	)

	return nil
}

func (b *memoryBus) Publish(
	e Event,
) error {
	if e.Name == "" {
		return ErrEventNameEmpty
	}

	b.mu.RLock()
	handlers, ok := b.handlers[e.Name]
	b.mu.RUnlock()

	if !ok {
		return nil
	}

	for _, handler := range handlers {
		go execute(handler, e)
	}

	return nil
}

// execute runs a handler in a panic-safe manner.
func execute(
	h Handler,
	e Event,
) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("handler panicked: %v, event: %+v", r, e)
		}
	}()
	_ = h.Handle(e)
}