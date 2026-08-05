package context

import (
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()

	if ctx == nil {
		t.Fatal("NewContext returned nil")
	}

	// Test Registry access
	reg := ctx.Registry()
	if reg == nil {
		t.Error("Context.Registry() returned nil")
	}
	// Further checks could involve adding an object to the registry and retrieving it
	// For now, just checking if it's not nil is sufficient for context initialization.

	// Test EventBus access
	eb := ctx.EventBus()
	if eb == nil {
		t.Error("Context.EventBus() returned nil")
	}
	// Similar to Registry, further checks would involve publishing/subscribing
	// For now, just checking if it's not nil is sufficient.
}