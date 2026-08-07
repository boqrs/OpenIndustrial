package event

import (
	"sync"
	"testing"
	"time"
)

// testHandler is a mock handler for testing purposes.
type testHandler struct {
	handledEvent *Event
	wg           sync.WaitGroup
	mu           sync.Mutex
}

func (h *testHandler) Handle(event Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handledEvent = &event
	h.wg.Done()
	return nil
}

func TestMemoryBusPublish(t *testing.T) {
	bus := NewMemoryBus()

	handler := &testHandler{}
	handler.wg.Add(1)

	eventName := "test.event"

	err := bus.Subscribe(eventName, handler)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	eventToPublish := Event{
		Name: eventName,
		Data: "test data",
	}

	err = bus.Publish(eventToPublish)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// Wait for the handler to process the event, with a timeout
	waitChan := make(chan struct{})
	go func() {
		handler.wg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for handler to be called")
	}

	handler.mu.Lock()
	defer handler.mu.Unlock()

	if handler.handledEvent == nil {
		t.Fatal("handler was not called")
	}

	if handler.handledEvent.Name != eventName {
		t.Errorf("handler received wrong event name: got %v, want %v", handler.handledEvent.Name, eventName)
	}

	if handler.handledEvent.Data != "test data" {
		t.Errorf("handler received wrong event data: got %v, want %v", handler.handledEvent.Data, "test data")
	}
}

func TestMemoryBusSubscribeValidation(t *testing.T) {
	bus := NewMemoryBus()
	handler := &testHandler{}

	if err := bus.Subscribe("", handler); err != ErrEventNameEmpty {
		t.Errorf("expected ErrEventNameEmpty for empty event name, got %v", err)
	}

	if err := bus.Subscribe("some.event", nil); err != ErrHandlerNil {
		t.Errorf("expected ErrHandlerNil for nil handler, got %v", err)
	}
}

func TestMemoryBusPublishValidation(t *testing.T) {
	bus := NewMemoryBus()
	if err := bus.Publish(Event{Name: ""}); err != ErrEventNameEmpty {
		t.Errorf("expected ErrEventNameEmpty for empty event name, got %v", err)
	}
}

func TestMemoryBusPanicRecovery(t *testing.T) {
	bus := NewMemoryBus()
	
	panicHandler := &panicHandler{}
	
	err := bus.Subscribe("panic.event", panicHandler)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// This should not crash the test process
	err = bus.Publish(Event{Name: "panic.event"})
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// Give the goroutine a moment to run and potentially panic
	time.Sleep(100 * time.Millisecond)
}

type panicHandler struct{}

func (h *panicHandler) Handle(event Event) error {
	panic("handler intentionally panicked")
}