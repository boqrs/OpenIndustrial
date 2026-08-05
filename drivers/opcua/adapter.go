package opcua

import (
	"context"
)

type Adapter interface {
	Connect(ctx context.Context, cfg ConnectionConfig) error
	Disconnect(ctx context.Context) error
	ReadNodes(ctx context.Context, mappings []NodeMapping) ([]Sample, error)
	WriteNode(ctx context.Context, sample Sample) error
	Subscribe(ctx context.Context, subCfg SubscriptionConfig, mappings []NodeMapping, sampleCh chan<- Sample) error
	State() State
	IsConnected() bool
}