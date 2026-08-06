package iec104

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Driver is the main entry point for the IEC 104 module.
type Driver struct {
	config         *ConnectionConfig
	adapter        Adapter
	cache          *Cache
	infoObjectChan chan InformationObject
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	ioaToID        map[uint32]string
	idToPoint      map[string]PointMapping
}

// NewDriver creates a new IEC 104 driver and validates the configuration.
func NewDriver(cfg *ConnectionConfig, points []PointMapping) (*Driver, error) {
	adapter := NewTCPAdapter()
	cache := NewCache()
	infoObjectChan := make(chan InformationObject, 100)

	ioaToID := make(map[uint32]string, len(points))
	idToPoint := make(map[string]PointMapping, len(points))
	seenIOAs := make(map[uint32]struct{}, len(points))
	seenIDs := make(map[string]struct{}, len(points))

	for _, p := range points {
		if _, exists := seenIDs[p.ID]; exists {
			return nil, fmt.Errorf("duplicate point ID found in configuration: %s", p.ID)
		}
		if _, exists := seenIOAs[p.IOA]; exists {
			return nil, fmt.Errorf("duplicate IOA found in configuration: %d", p.IOA)
		}
		seenIDs[p.ID] = struct{}{}
		seenIOAs[p.IOA] = struct{}{}
		ioaToID[p.IOA] = p.ID
		idToPoint[p.ID] = p
	}

	return &Driver{
		config:         cfg,
		adapter:        adapter,
		cache:          cache,
		infoObjectChan: infoObjectChan,
		ioaToID:        ioaToID,
		idToPoint:      idToPoint,
	}, nil
}

// Start connects the adapter and starts the subscription and processing loops.
func (d *Driver) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel

	log.Println("Starting IEC 104 driver...")

	if err := d.adapter.Connect(ctx, *d.config); err != nil {
		log.Printf("Failed to connect IEC 104 adapter: %v", err)
		cancel()
		return err
	}
	log.Println("IEC 104 adapter connected successfully.")

	go func() {
		err := d.adapter.Subscribe(ctx, d.infoObjectChan)
		if err != nil && err != context.Canceled {
			log.Printf("IEC 104 subscription failed: %v", err)
		}
	}()

	d.wg.Add(1)
	go d.processInfoObjects(ctx)

	log.Println("IEC 104 driver started successfully.")
	return nil
}

// Stop stops the subscription and disconnects the adapter gracefully.
func (d *Driver) Stop() error {
	log.Println("Stopping IEC 104 driver...")
	if d.cancel != nil {
		d.cancel()
	}

	d.wg.Wait()

	if d.adapter != nil && d.adapter.IsConnected() {
		if err := d.adapter.Disconnect(context.Background()); err != nil {
			return err
		}
	}
	log.Println("IEC 104 driver stopped.")
	return nil
}

// processInfoObjects runs a loop to read InformationObjects and update the cache with Samples.
func (d *Driver) processInfoObjects(ctx context.Context) {
	defer d.wg.Done()
	for {
		select {
		case obj := <-d.infoObjectChan:
			if pointID, ok := d.ioaToID[obj.IOA]; ok {
				sample := Sample{
					PointID:   pointID,
					Value:     obj.Value,
					Timestamp: obj.Timestamp,
					Source:    "iec104",
				}
				if sample.Timestamp.IsZero() {
					sample.Timestamp = time.Now()
				}
				d.cache.Set(pointID, sample)
			} else {
				log.Printf("Warning: received data for unmapped IOA: %d", obj.IOA)
			}
		case <-ctx.Done():
			log.Println("Information object processing loop stopped.")
			return
		}
	}
}

// Read returns the latest value for a specific point ID from the cache.
func (d *Driver) Read(pointID string) (Sample, error) {
	sample, found := d.cache.Get(pointID)
	if !found {
		return Sample{}, fmt.Errorf("point ID '%s' not found in cache", pointID)
	}
	return sample, nil
}

// Write sends a command to a specific point using a fast map lookup.
func (d *Driver) Write(pointID string, value interface{}) error {
	targetPoint, found := d.idToPoint[pointID]
	if !found {
		return fmt.Errorf("point ID '%s' not found in mapping", pointID)
	}

	return d.adapter.Write(context.Background(), targetPoint, value)
}