//go:build !linux

package profinet

import (
	"context"
	"log"
)

// pnetAdapterStub is a placeholder implementation for non-Linux systems.
// It ensures the project can be compiled on any platform.
type pnetAdapterStub struct{}

// NewPnetAdapter returns a stub adapter for non-Linux systems.
func NewPnetAdapter() Adapter {
	log.Println("PROFINET is only supported on Linux. Using a stub adapter.")
	return &pnetAdapterStub{}
}

func (a *pnetAdapterStub) Connect(ctx context.Context, cfg ConnectionConfig) error {
	return ErrUnsupportedOnPlatform
}

func (a *pnetAdapterStub) Disconnect(ctx context.Context) error {
	return nil
}

func (a *pnetAdapterStub) IsConnected() bool {
	return false
}

func (a *pnetAdapterStub) Discover(ctx context.Context) ([]DeviceInfo, error) {
	return nil, ErrUnsupportedOnPlatform
}

func (a *pnetAdapterStub) ConnectDevice(ctx context.Context, device DeviceInfo) error {
	return ErrUnsupportedOnPlatform
}

func (a *pnetAdapterStub) ReadInputs(ctx context.Context, deviceID string) ([]byte, error) {
	return nil, ErrUnsupportedOnPlatform
}

func (a *pnetAdapterStub) WriteOutputs(ctx context.Context, deviceID string, data []byte) error {
	return ErrUnsupportedOnPlatform
}

func (a *pnetAdapterStub) ReadRecord(ctx context.Context, req RecordRequest) (RecordResponse, error) {
	return RecordResponse{}, ErrUnsupportedOnPlatform
}

func (a *pnetAdapterStub) Subscribe(ctx context.Context, ch chan<- Sample) error {
	return ErrUnsupportedOnPlatform
}