package service

import (
	"context"

	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/lifecycle"
)

// Service defines the interface for a runtime service.
// All services must implement the Lifecycle interface to be managed by the runtime.
type Service interface {
	lifecycle.Lifecycle
	// Name returns the unique name of the service.
	Name() string
}

// BaseService provides a basic implementation for the Service interface,
// allowing embedding in custom service structs to reduce boilerplate.
type BaseService struct {
	ServiceName string
}

// NewBaseService creates a new BaseService instance.
func NewBaseService(name string) *BaseService {
	return &BaseService{
		ServiceName: name,
	}
}

// Name returns the name of the service.
func (bs *BaseService) Name() string {
	return bs.ServiceName
}

// Start is a no-op default implementation for BaseService.
// Custom services should override this method with their specific startup logic.
func (bs *BaseService) Start(ctx context.Context) error {
	// Default: do nothing, just return nil
	return nil
}

// Stop is a no-op default implementation for BaseService.
// Custom services should override this method with their specific shutdown logic.
func (bs *BaseService) Stop(ctx context.Context) error {
	// Default: do nothing, just return nil
	return nil
}