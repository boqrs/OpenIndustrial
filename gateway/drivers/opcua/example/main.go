package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway/drivers/opcua"
)

func main() {
	// 1. Configure the driver
	// We will connect to a public OPC UA test server and poll a single node.
	cfg := opcua.DriverConfig{
		Connection: opcua.ConnectionConfig{
			EndpointURL: "opc.tcp://opcuademo.sterfive.com:26543",
			AuthType:    "Anonymous",
			Timeout:     5 * time.Second,
		},
		Poll: opcua.PollConfig{
			Interval: 5 * time.Second,
		},
		Mappings: []opcua.NodeMapping{
			{
				ID:     "server_current_time",
				NodeID: "ns=1;i=1013", // This node provides the server's current time
			},
		},
		CollectionMode: "poll", // Use polling mode
	}

	// 2. Create a new driver instance
	// We use the NoOpPublisher for this example, which just logs the data.
	driver, err := opcua.NewDriver("opcua-test-driver", cfg, opcua.NewNoOpPublisher())
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	// 3. Start the driver
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver: %v", err)
	}

	// Set up a channel to listen for OS signals (like Ctrl+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Ticker to periodically read data from the driver's cache
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("Driver is running. Press Ctrl+C to exit.")

	// 4. Main loop
	// We'll periodically get the value from the cache and listen for the shutdown signal.
Loop:
	for {
		select {
		case <-ticker.C:
			// Get the value from the cache
			sample, ok := driver.GetValue("server_current_time")
			if ok {
				log.Printf("--- Reading from cache ---")
				log.Printf("ID: %s, Value: %v, Timestamp: %s, Quality: %s",
					sample.ID, sample.Value, sample.Timestamp.Format(time.RFC3339), sample.Quality)
				log.Printf("--------------------------")
			} else {
				log.Println("Value not yet available in cache.")
			}
		case <-sigCh:
			log.Println("Shutdown signal received.")
			break Loop
		}
	}

	// 5. Stop the driver gracefully
	log.Println("Stopping driver...")
	if err := driver.Stop(); err != nil {
		log.Fatalf("Failed to stop driver gracefully: %v", err)
	}

	log.Println("Driver stopped. Exiting.")
}