package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/drivers/simulator"
	"github.com/OpenGongChang/OpenIndustrial/gateway"
	"github.com/OpenGongChang/OpenIndustrial/runtime/driver"
)

func main() {
	log.Println("Starting OpenIndustrial Gateway Example...")

	// 1. Define configuration
	// In a real application, this would be loaded from a file (e.g., config.yaml)
	cfg := &gateway.GatewayConfig{
		ID:   "gateway-01",
		Name: "Main Gateway",
		Drivers: []driver.Config{
			{
				ID:   "multi-sim-1",
				Type: simulator.DriverType,
				Settings: map[string]interface{}{
					"devices": []map[string]interface{}{
						{
							"deviceId": "hvac-01",
							"points": []map[string]interface{}{
								{
									"pointId":   "temperature",
									"valueType": "float",
									"min":       18.0,
									"max":       25.0,
									"interval":  "5s",
								},
								{
									"pointId":   "running_status",
									"valueType": "bool",
									"interval":  "10s",
								},
							},
						},
						{
							"deviceId": "power-meter-01",
							"points": []map[string]interface{}{
								{
									"pointId":   "voltage",
									"valueType": "sine",
									"base":      220.0,
									"amplitude": 5.0,
									"period":    "60s",
									"interval":  "2s",
								},
							},
						},
					},
				},
			},
		},
	}

	// 2. Create a new gateway instance
	gw := gateway.NewGateway(cfg)

	// 3. Register driver factories
	// This is how the gateway knows how to create a driver of a specific type.
	err := gw.Registry().Register(simulator.DriverType, func(c driver.Config) (driver.Driver, error) {
		return simulator.NewDriver(c)
	})
	if err != nil {
		log.Fatalf("Failed to register simulator driver: %v", err)
	}

	// 4. Start the gateway
	// This will initialize and start all configured drivers.
	if err := gw.Start(); err != nil {
		log.Fatalf("Failed to start gateway: %v", err)
	}

	// 5. Set up a listener for data from the collector
	// This simulates a component (like a publisher) that consumes the data.
	go func() {
		subscriber := gw.Collector().Subscribe()
		log.Println("Data subscriber started. Waiting for events from the collector...")
		for event := range subscriber {
			log.Printf("[SUBSCRIBER] Received Event: DriverID=%s, DeviceID=%s, PointID=%s, Value=%.2f",
				event.DriverID, event.DeviceID, event.PointID, event.Value)
		}
		log.Println("Data subscriber stopped.")
	}()


	// 6. Wait for shutdown signal
	// This keeps the main function alive and allows for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gateway...")

	// 7. Stop the gateway
	if err := gw.Stop(); err != nil {
		log.Printf("Error during gateway shutdown: %v", err)
	}

	log.Println("Gateway has been shut down gracefully.")
	// Give a moment for logs to be flushed
	time.Sleep(1 * time.Second)
}