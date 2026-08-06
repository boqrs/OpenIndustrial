package opcua

import (
	"context"
	"fmt"
	"log"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type opcuaClientAdapter struct {
	BaseAdapter
	client *opcua.Client
}

func NewOpcuaClientAdapter() Adapter {
	return &opcuaClientAdapter{}
}

func (a *opcuaClientAdapter) Connect(ctx context.Context, cfg ConnectionConfig) error {
	if a.IsConnected() {
		return ErrAlreadyConnected
	}
	a.SetState(StateConnecting)

	if err := cfg.Validate(); err != nil {
		a.SetState(StateDisconnected)
		return err
	}

	opts := []opcua.Option{
		opcua.SecurityPolicy(cfg.SecurityPolicy),
		opcua.SecurityModeString(cfg.SecurityMode),
		opcua.AuthUsername(cfg.Username, cfg.Password),
		opcua.CertificateFile(cfg.CertPath),
		opcua.PrivateKeyFile(cfg.KeyPath),
		opcua.RequestTimeout(cfg.Timeout),
	}

	client, err := opcua.NewClient(cfg.EndpointURL, opts...)
	if err != nil {
		a.SetState(StateDisconnected)
		return fmt.Errorf("failed to create OPC UA client: %w", err)
	}

	if err := client.Connect(ctx); err != nil {
		a.SetState(StateDisconnected)
		return fmt.Errorf("failed to connect to OPC UA server: %w", err)
	}

	a.client = client
	a.SetState(StateConnected)
	log.Printf("Successfully connected to OPC UA server at %s", cfg.EndpointURL)
	return nil
}

func (a *opcuaClientAdapter) Disconnect(ctx context.Context) error {
	if !a.IsConnected() {
		return nil
	}
	if err := a.client.Close(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from OPC UA server: %w", err)
	}
	a.SetState(StateDisconnected)
	log.Println("Disconnected from OPC UA server")
	return nil
}

func (a *opcuaClientAdapter) ReadNodes(ctx context.Context, mappings []NodeMapping) ([]Sample, error) {
	if !a.IsConnected() {
		return nil, ErrNotConnected
	}

	nodesToRead := make([]*ua.ReadValueID, 0, len(mappings))
	validMappings := make([]NodeMapping, 0, len(mappings))
	for _, nm := range mappings {
		nodeID, err := ua.ParseNodeID(nm.NodeID)
		if err != nil {
			log.Printf("Skipping invalid NodeID '%s' for reading: %v", nm.NodeID, err)
			continue
		}
		nodesToRead = append(nodesToRead, &ua.ReadValueID{NodeID: nodeID, AttributeID: ua.AttributeIDValue})
		validMappings = append(validMappings, nm)
	}

	if len(nodesToRead) == 0 {
		return []Sample{}, nil
	}

	req := &ua.ReadRequest{
		MaxAge:             2000,
		NodesToRead:        nodesToRead,
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	resp, err := a.client.Read(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("OPC UA read failed: %w", err)
	}

	samples := make([]Sample, 0, len(resp.Results))
	for i, result := range resp.Results {
		if result.Status != ua.StatusGood {
			log.Printf("Failed to read node %s: %s", validMappings[i].NodeID, result.Status)
		}
		samples = append(samples, Sample{
			ID:        validMappings[i].ID,
			NodeID:    validMappings[i].NodeID,
			Value:     result.Value.Value(),
			Timestamp: result.SourceTimestamp,
			Quality:   ToQuality(result.Status),
		})
	}

	return samples, nil
}

func (a *opcuaClientAdapter) WriteNode(ctx context.Context, sample Sample) error {
	if !a.IsConnected() {
		return ErrNotConnected
	}

	nodeID, err := ua.ParseNodeID(sample.NodeID)
	if err != nil {
		return fmt.Errorf("invalid NodeID '%s' for writing: %w", sample.NodeID, err)
	}

	variant, err := ua.NewVariant(sample.Value)
	if err != nil {
		return fmt.Errorf("value for node %s cannot be converted to variant: %w", sample.NodeID, err)
	}

	req := &ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{
			{
				NodeID:      nodeID,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					Value: variant,
				},
			},
		},
	}

	resp, err := a.client.Write(ctx, req)
	if err != nil {
		return fmt.Errorf("OPC UA write failed: %w", err)
	}

	if resp.Results[0] != ua.StatusGood {
		return fmt.Errorf("failed to write to node %s: %s", sample.NodeID, resp.Results[0])
	}

	return nil
}