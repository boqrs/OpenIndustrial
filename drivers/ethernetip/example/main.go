package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/drivers/ethernetip"
)

func main() {
	log.Println("Starting EtherNet/IP driver example...")

	// 1. Define the connection configuration for the PLC.
	// Replace with your PLC's actual IP address.
	connConfig := &ethernetip.ConnectionConfig{
		Host:    "192.168.1.10", // IMPORTANT: Change to your PLC's IP
		Port:    44818,
		Slot:    0,
		Timeout: 3 * time.Second,
		Mode:    ethernetip.ModeExplicit,
	}

	// 2. Define the points (tags) we want to read from the PLC.
	pointsToRead := []ethernetip.PointMapping{
		{
			ID:   "motor_speed",
			Type: ethernetip.PointTypePLCTag,
			Tag:  "Motor.Speed", // Example tag
		},
		{
			ID:   "tank_level",
			Type: ethernetip.PointTypePLCTag,
			Tag:  "Tank.Level", // Example tag
		},
		{
			ID:   "valve_status",
			Type: ethernetip.PointTypePLCTag,
			Tag:  "Valve.Status", // Example tag
		},
	}

	// 3. Create a new driver instance.
	driver := ethernetip.NewDriver(connConfig, pointsToRead)

	// 4. Start the driver. This will connect to the PLC and start polling.
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start EtherNet/IP driver: %v", err)
	}
	defer driver.Stop()

	// 5. Start a goroutine to consume samples from the driver.
	// The Samples() channel provides a stream of data that has already
	// passed through the Optimizer, so we only get values when they change.
	go func() {
		log.Println("Starting sample consumer...")
		sampleChan := driver.Samples()
		for sample := range sampleChan {
			log.Printf("--> Received Sample: PointID=[%s], Value=[%v], Quality=[%s], Timestamp=[%s]",
				sample.PointID,
				sample.Value,
				sample.Quality,
				sample.Timestamp.Format(time.RFC3339),
			)
		}
		log.Println("Sample consumer finished.")
	}()

	// 6. Wait for a shutdown signal to gracefully stop the driver.
	log.Println("Example is running. Press Ctrl+C to stop.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown signal received. Stopping driver...")
}