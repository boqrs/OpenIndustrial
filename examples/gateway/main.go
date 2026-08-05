package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/drivers/simulator"
	"github.com/OpenGongChang/OpenIndustrial/gateway"
	_ "github.com/OpenGongChang/OpenIndustrial/gateway/publishers/mqtt" // Anonymous import to run init()
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
		Publishers: []gateway.PublisherConfig{
			{
				ID:   "mqtt-publisher-1",
				Type: "mqtt",
				Settings: map[string]interface{}{
					"broker":   "tcp://localhost:1883",
					"clientId": "gateway-01", // Should be unique
					"qos":      1,
				},
			},
		},
	}

	// 2. Create a new gateway instance
	gw, err := gateway.NewGateway(cfg)
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	// 3. Register driver factories
	// This is how the gateway knows how to create a driver of a specific type.
	err = gw.Registry().Register(simulator.DriverType, func(c driver.Config) (driver.Driver, error) {
		return simulator.NewDriver(c)
	})
	if err != nil {
		log.Fatalf("Failed to register simulator driver: %v", err)
	}

	// 4. Start the gateway
	// This will initialize and start all configured drivers and publishers.
	if err := gw.Start(); err != nil {
		log.Fatalf("Failed to start gateway: %v", err)
	}

	// 5. Wait for shutdown signal
	// This keeps the main function alive and allows for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gateway...")

	// 6. Stop the gateway
	if err := gw.Stop(); err != nil {
		log.Printf("Error during gateway shutdown: %v", err)
	}

	log.Println("Gateway has been shut down gracefully.")
	// Give a moment for logs to be flushed
	time.Sleep(1 * time.Second)
}