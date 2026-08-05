package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

// CustomService is a mock service that embeds BaseService and overrides Start/Stop.
type CustomService struct {
	*BaseService
	startCalled bool
	stopCalled  bool
	startErr    error
	stopErr     error
	startDelay  time.Duration
	stopDelay   time.Duration
}

// NewCustomService creates a new CustomService.
func NewCustomService(name string) *CustomService {
	return &CustomService{
		BaseService: NewBaseService(name),
	}
}

// Start overrides the BaseService's Start method.
func (cs *CustomService) Start(ctx context.Context) error {
	cs.startCalled = true
	if cs.startDelay > 0 {
		select {
		case <-time.After(cs.startDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return cs.startErr
}

// Stop overrides the BaseService's Stop method.
func (cs *CustomService) Stop(ctx context.Context) error {
	cs.stopCalled = true
	if cs.stopDelay > 0 {
		select {
		case <-time.After(cs.stopDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return cs.stopErr
}

func TestBaseService(t *testing.T) {
	bs := NewBaseService("TestBaseService")

	if bs.Name() != "TestBaseService" {
		t.Errorf("Expected name 'TestBaseService', got '%s'", bs.Name())
	}

	ctx := context.Background()
	if err := bs.Start(ctx); err != nil {
		t.Errorf("BaseService Start should not return error, got %v", err)
	}
	if err := bs.Stop(ctx); err != nil {
		t.Errorf("BaseService Stop should not return error, got %v", err)
	}
}

func TestCustomService(t *testing.T) {
	cs := NewCustomService("TestCustomService")

	if cs.Name() != "TestCustomService" {
		t.Errorf("Expected name 'TestCustomService', got '%s'", cs.Name())
	}

	ctx := context.Background()

	// Test successful start/stop
	err := cs.Start(ctx)
	if err != nil {
		t.Errorf("Expected no error on Start, got %v", err)
	}
	if !cs.startCalled {
		t.Error("CustomService Start method was not called")
	}

	err = cs.Stop(ctx)
	if err != nil {
		t.Errorf("Expected no error on Stop, got %v", err)
	}
	if !cs.stopCalled {
		t.Error("CustomService Stop method was not called")
	}

	// Test start with error
	cs = NewCustomService("TestCustomServiceWithError")
	cs.startErr = errors.New("custom start failed")
	err = cs.Start(ctx)
	if err == nil {
		t.Error("Expected error on Start, got nil")
	}
	if err.Error() != "custom start failed" {
		t.Errorf("Expected specific error, got %v", err)
	}

	// Test stop with context cancellation
	cs = NewCustomService("TestCustomServiceWithCancel")
	cs.stopDelay = time.Second
	cancelCtx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err = cs.Stop(cancelCtx)
	if err == nil {
		t.Error("Expected context cancellation error on Stop, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}