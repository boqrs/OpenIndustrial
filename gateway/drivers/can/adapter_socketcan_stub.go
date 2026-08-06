//go:build !linux

package can

import (
	"context"
	"errors"
)

// NewSocketCANAdapter is a stub for non-linux systems
func NewSocketCANAdapter() Adapter {
	return &socketcanAdapterStub{}
}

type socketcanAdapterStub struct{}

func (a *socketcanAdapterStub) Connect(ctx context.Context, cfg ConnectionConfig) error {
	return errors.New("SocketCAN is not supported on this platform")
}

func (a *socketcanAdapterStub) Disconnect(ctx context.Context) error {
	return nil
}

func (a *socketcanAdapterStub) IsConnected() bool {
	return false
}

func (a *socketcanAdapterStub) Receive() <-chan Frame {
	// Return a nil channel, which will block forever and never yield any frames.
	return nil
}

func (a *socketcanAdapterStub) Send(ctx context.Context, frame Frame) error {
	return errors.New("SocketCAN is not supported on this platform")
}