package lifecycle

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockLifecycleComponent is a mock implementation of the Lifecycle interface for testing.
type MockLifecycleComponent struct {
	Name        string
	StartCalled bool
	StopCalled  bool
	StartError  error
	StopError   error
	StartDelay  time.Duration
	StopDelay   time.Duration
}

// NewMockLifecycleComponent creates a new mock component.
func NewMockLifecycleComponent(name string) *MockLifecycleComponent {
	return &MockLifecycleComponent{Name: name}
}

// Start simulates the start operation of a component.
func (m *MockLifecycleComponent) Start(ctx context.Context) error {
	m.StartCalled = true
	if m.StartDelay > 0 {
		select {
		case <-time.After(m.StartDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return m.StartError
}

// Stop simulates the stop operation of a component.
func (m *MockLifecycleComponent) Stop(ctx context.Context) error {
	m.StopCalled = true
	if m.StopDelay > 0 {
		select {
		case <-time.After(m.StopDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return m.StopError
}

func TestLifecycleStart(t *testing.T) {
	ctx := context.Background()
	comp := NewMockLifecycleComponent("TestComp")

	// Test successful start
	err := comp.Start(ctx)
	if err != nil {
		t.Errorf("Expected no error on Start, got %v", err)
	}
	if !comp.StartCalled {
		t.Error("Start method was not called")
	}

	// Test start with error
	comp = NewMockLifecycleComponent("TestCompWithError")
	comp.StartError = errors.New("failed to initialize")
	err = comp.Start(ctx)
	if err == nil {
		t.Error("Expected error on Start, got nil")
	}
	if err.Error() != "failed to initialize" {
		t.Errorf("Expected specific error, got %v", err)
	}
}

func TestLifecycleStop(t *testing.T) {
	ctx := context.Background()
	comp := NewMockLifecycleComponent("TestComp")

	// Test successful stop
	err := comp.Stop(ctx)
	if err != nil {
		t.Errorf("Expected no error on Stop, got %v", err)
	}
	if !comp.StopCalled {
		t.Error("Stop method was not called")
	}

	// Test stop with error
	comp = NewMockLifecycleComponent("TestCompWithError")
	comp.StopError = errors.New("failed to shutdown")
	err = comp.Stop(ctx)
	if err == nil {
		t.Error("Expected error on Stop, got nil")
	}
	if err.Error() != "failed to shutdown" {
		t.Errorf("Expected specific error, got %v", err)
	}
}

func TestLifecycleStartWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	comp := NewMockLifecycleComponent("TestCompWithDelay")
	comp.StartDelay = time.Second // Simulate a long startup

	go func() {
		time.Sleep(100 * time.Millisecond) // Cancel context before delay finishes
		cancel()
	}()

	err := comp.Start(ctx)
	if err == nil {
		t.Error("Expected context cancellation error on Start, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
	if !comp.StartCalled {
		t.Error("Start method was not called")
	}
}

func TestLifecycleStopWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	comp := NewMockLifecycleComponent("TestCompWithDelay")
	comp.StopDelay = time.Second // Simulate a long shutdown

	go func() {
		time.Sleep(100 * time.Millisecond) // Cancel context before delay finishes
		cancel()
	}()

	err := comp.Stop(ctx)
	if err == nil {
		t.Error("Expected context cancellation error on Stop, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
	if !comp.StopCalled {
		t.Error("Stop method was not called")
	}
}