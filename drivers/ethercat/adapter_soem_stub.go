//go:build !linux
// +build !linux

package ethercat

import (
	"context"
	"fmt"
	"log"
)

// soemAdapter is the stub implementation for non-cgo builds.
type soemAdapter struct{}

// NewSOEMAdapter returns a new stub adapter.
func NewSOEMAdapter() Adapter {
	log.Println("WARNING: You are using the CGO-less stub for the EtherCAT driver. No real hardware communication will occur.")
	return &soemAdapter{}
}

func (a *soemAdapter) Connect(ctx context.Context, cfg ConnectionConfig) error {
	return fmt.Errorf("ethercat: Connect not implemented in non-cgo build")
}

func (a *soemAdapter) Disconnect(ctx context.Context) error {
	return fmt.Errorf("ethercat: Disconnect not implemented in non-cgo build")
}

func (a *soemAdapter) IsConnected() bool {
	return false
}

func (a *soemAdapter) ScanSlaves(ctx context.Context) ([]SlaveInfo, error) {
	return nil, fmt.Errorf("ethercat: ScanSlaves not implemented in non-cgo build")
}

func (a *soemAdapter) ConfigurePDOs(ctx context.Context, mappings []PDOMapping) error {
	return fmt.Errorf("ethercat: ConfigurePDOs not implemented in non-cgo build")
}

func (a *soemAdapter) ReadPDOs(ctx context.Context) (map[uint16][]byte, error) {
	return nil, fmt.Errorf("ethercat: ReadPDOs not implemented in non-cgo build")
}

func (a *soemAdapter) WritePDOs(ctx context.Context, data map[uint16][]byte) error {
	return fmt.Errorf("ethercat: WritePDOs not implemented in non-cgo build")
}

func (a *soemAdapter) ReadSDO(ctx context.Context, req SDORequest) (SDOResponse, error) {
	return SDOResponse{}, fmt.Errorf("ethercat: ReadSDO not implemented in non-cgo build")
}

func (a *soemAdapter) WriteSDO(ctx context.Context, req SDORequest) error {
	return fmt.Errorf("ethercat: WriteSDO not implemented in non-cgo build")
}