package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/drivers/modbus"
	"github.com/OpenGongChang/OpenIndustrial/drivers/modbus/internal/codec" // Import internal/codec for decoding
	runtimeContext "github.com/OpenGongChang/OpenIndustrial/runtime/context"
	"github.com/OpenGongChang/OpenIndustrial/runtime/eventbus"
)

func main() {
	log.Println("Starting Modbus example driver...")

	// 1. Configuration
	// Define your Modbus TCP device details and the points you want to read.
	// This example assumes a Modbus TCP slave is running at localhost:502.
	cfg := modbus.Config{
		Name: "Modbus-TCP-Driver-1",
		Connection: modbus.ConnectionConfig{
			Mode:        modbus.ModeTCP,
			Address:     "127.0.0.1:502",
			Timeout:     5 * time.Second,
			SlaveID:     1,
		},
		Poll: modbus.PollConfig{
			Interval: 1 * time.Second,
		},
		NodeMappings: []modbus.NodeMapping{
			{
				PointID:    "temperature_coil",
				Register:   modbus.RegisterCoil,
				Address:    0,
				DataType:   codec.DataTypeBool,
				Length:     1,
				Writable:   true,
				Description: "Example Coil for Temperature",
			},
			{
				PointID:    "humidity_holding",
				Register:   modbus.RegisterHoldingRegister,
				Address:    1,
				DataType:   codec.DataTypeUint16,
				Length:     1,
				ByteOrder:  codec.BigEndian,
				WordOrder:  codec.LowWordFirst,
				Writable:   false,
				Description: "Example Holding Register for Humidity",
			},
			{
				PointID:    "pressure_input",
				Register:   modbus.RegisterInputRegister,
				Address:    2,
				DataType:   codec.DataTypeUint16,
				Length:     1,
				ByteOrder:  codec.BigEndian,
				WordOrder:  codec.LowWordFirst,
				Writable:   false,
				Description: "Example Input Register for Pressure",
			},
		},
	}

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid Modbus configuration: %v", err)
	}

	// 2. Driver Initialization
	driver, err := modbus.NewDriver(cfg)
	if err != nil {
		log.Fatalf("Failed to create Modbus driver: %v", err)
	}

	// 3. Runtime Context (minimal for example)
	// In a real application, this would be provided by the OpenIndustrial runtime.
	eb := eventbus.NewEventBus()
	ctx := runtimeContext.NewContext()

	// Subscribe to driver events
	eventCh := make(chan eventbus.Event, 10)
	eb.Subscribe(driver.GetID(), func(event eventbus.Event) {
		eventCh <- event
	})
	log.Printf("Subscribed to Modbus driver events on topic: %s", driver.GetID())

	// Goroutine to consume events
	go func() {
		for {
			select {
			case event := <-eventCh:
				sample, ok := event.Data.(modbus.Sample)
				if !ok {
					log.Printf("Received non-Modbus sample event: %+v", event)
					continue
				}
				log.Printf("Received Modbus Sample: PointID=%s, Value=%v, Timestamp=%s, Quality=%s",
					sample.PointID, sample.Value, sample.Timestamp.Format(time.RFC3339), sample.Quality)
			}
		}
	}()

	// Context for driver lifecycle management
	driverCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 4. Driver Lifecycle: Init and Start
	if err := driver.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize Modbus driver: %v", err)
	}
	log.Println("Modbus driver initialized.")

	if err := driver.Start(driverCtx); err != nil {
		log.Fatalf("Failed to start Modbus driver: %v", err)
	}
	log.Println("Modbus driver started. Polling for data...")

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down Modbus driver...")
	if err := driver.Stop(driverCtx); err != nil {
		log.Printf("Error stopping Modbus driver: %v", err)
	}
	log.Println("Modbus driver stopped.")
	log.Println("Modbus example driver finished.")
}