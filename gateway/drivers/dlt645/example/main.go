package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway/drivers/dlt645"
)

func main() {
	log.Println("--- DL/T 645 Driver Example ---")

	// 1. Configuration
	// IMPORTANT: Change this to your actual serial port address.
	// On Linux, it might be /dev/ttyUSB0. On Windows, COM3.
	serialPortAddress := "/dev/tty.usbserial-0001"

	config := &dlt645.ConnectionConfig{
		Type:     dlt645.Serial,
		Address:  serialPortAddress,
		Timeout:  5 * time.Second,
		BaudRate: 9600,
		DataBits: 8,
		StopBits: 1,
		Parity:   "E", // Even parity is common for DL/T 645
	}

	// Define the meters and data points you want to read.
	meters := []dlt645.Meter{
		{Address: "000000123456"}, // Use the actual 12-digit BCD address of your meter
	}

	points := []dlt645.PointMapping{
		{ID: "total_energy", DI: 0x00000000, DataType: "bcd", Scale: 0.01},
		// Add more points here if needed
		// {ID: "voltage_a", DI: 0x02010100, DataType: "bcd", Scale: 0.1},
	}

	// 2. Initialization
	driver, err := dlt645.NewDriver(config, meters, points)
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	// 3. Start the driver
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver: %v", err)
	}
	defer driver.Stop()

	log.Printf("Driver started. Polling meter(s) on port %s every %s.",
		config.Address, config.Timeout)
	log.Println("Press Ctrl+C to exit.")

	// Set up a channel for graceful shutdown.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// 4. Main loop to periodically read from cache
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			// Construct the point ID as used in the cache: {meter_address}.{point_id}
			pointID := fmt.Sprintf("%s.%s", meters[0].Address, points[0].ID)
			sample, err := driver.Read(pointID)
			if err != nil {
				log.Printf("Could not read '%s' from cache: %v", pointID, err)
			} else {
				log.Printf(">>> Latest data from cache: PointID=%s, Value=%v, Timestamp=%s",
					sample.PointID, sample.Value, sample.Timestamp.Format(time.RFC3339))
			}
		case <-stopChan:
			log.Println("Shutdown signal received.")
			break loop
		}
	}

	log.Println("Example finished.")
}