package opcua

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

func (a *opcuaClientAdapter) Subscribe(ctx context.Context, subCfg SubscriptionConfig, mappings []NodeMapping, sampleCh chan<- Sample) error {
	if !a.IsConnected() {
		return ErrNotConnected
	}

	notifyCh := make(chan *opcua.PublishNotificationData, 16)

	params := &opcua.SubscriptionParameters{
		Interval:          subCfg.PublishingInterval,
		LifetimeCount:     subCfg.LifetimeCount,
		MaxKeepAliveCount: subCfg.MaxKeepAliveCount,
	}

	subscription, err := a.client.Subscribe(ctx, params, notifyCh)
	if err != nil {
		return fmt.Errorf("OPC UA subscription failed: %w", err)
	}

	log.Printf("Created OPC UA subscription with ID: %d", subscription.SubscriptionID)

	clientHandleToMapping := make(map[uint32]NodeMapping)
	var clientHandle uint32 = 1

	for _, mapping := range mappings {
		nodeID, err := ua.ParseNodeID(mapping.NodeID)
		if err != nil {
			log.Printf("Skipping invalid NodeID '%s' for subscription: %v", mapping.NodeID, err)
			continue
		}

		req := opcua.NewMonitoredItemCreateRequestWithDefaults(nodeID, ua.AttributeIDValue, clientHandle)
		req.RequestedParameters.SamplingInterval = float64(subCfg.SamplingInterval / time.Millisecond)
		req.RequestedParameters.QueueSize = 10
		req.RequestedParameters.DiscardOldest = true

		clientHandleToMapping[clientHandle] = mapping
		clientHandle++

		// Using subscription.Monitor to create monitored items
		_, err = subscription.Monitor(ctx, ua.TimestampsToReturnBoth, req)
		if err != nil {
			log.Printf("Failed to monitor node %s: %v", mapping.NodeID, err)
			// We don't return here, just log the error and continue with other nodes.
		}
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Subscription context cancelled, stopping notification listener.")
				// The subscription is automatically cleaned up by the gopcua library when the context is cancelled.
				return
			case notification := <-notifyCh:
				if notification == nil {
					continue
				}
				if notification.Error != nil {
					log.Printf("Received subscription error: %s", notification.Error)
					continue
				}

				if dataChange, ok := notification.Value.(*ua.DataChangeNotification); ok {
					for _, item := range dataChange.MonitoredItems {
						if item.Value == nil || item.Value.Value == nil {
							continue
						}
						if mapping, ok := clientHandleToMapping[item.ClientHandle]; ok {
							sample := Sample{
								ID:        mapping.ID,
								NodeID:    mapping.NodeID,
								Value:     item.Value.Value.Value(),
								Timestamp: item.Value.SourceTimestamp,
								Quality:   ToQuality(item.Value.Status),
							}
							select {
							case sampleCh <- sample:
							case <-ctx.Done():
								return
							default:
								log.Printf("Sample channel is full, dropping sample for node %s", mapping.NodeID)
							}
						}
					}
				}
			}
		}
	}()

	return nil
}