package context

import (
	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/eventbus"
	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/registry"
)

// Context defines the core runtime context, providing access to essential services.
type Context interface {
	Registry() registry.Registry
	EventBus() eventbus.EventBus
	// Add other core services here as they are developed
}

// runtimeContext is the concrete implementation of the Context interface.
type runtimeContext struct {
	reg registry.Registry
	eb  eventbus.EventBus
}

// NewContext creates a new runtime Context with initialized services.
func NewContext() Context {
	return &runtimeContext{
		reg: registry.NewRegistry(),
		eb:  eventbus.NewEventBus(),
	}
}

// Registry returns the global object registry.
func (rc *runtimeContext) Registry() registry.Registry {
	return rc.reg
}

// EventBus returns the global event bus.
func (rc *runtimeContext) EventBus() eventbus.EventBus {
	return rc.eb
}