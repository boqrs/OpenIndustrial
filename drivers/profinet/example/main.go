package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/drivers/profinet"
)

func main() {
	log.Println("Starting PROFINET driver example...")

	// 1. Define the connection configuration for the PROFINET controller.
	// On Linux, "eth0" should be replaced with your actual network interface name.
	connConfig := &profinet.ConnectionConfig{
		Interface: "eth0",
		LocalIP:   "192.168.1.1", // The IP of our controller
		CycleTime: 4 * time.Millisecond,
		Timeout:   1 * time.Second,
	}

	// 2. Define the devices we want to communicate with.
	// This configuration would typically be loaded from a file.
	deviceConfigs := []profinet.DeviceConfig{
		{
			Name:        "ET200SP-Input-Module",
			IP:          "192.168.1.10",
			StationName: "et200sp-im",
			GSDML:       "GSDML-V2.3-Siemens-ET200SP-IM155-6PN-ST.xml",
			Modules: []profinet.ModuleConfig{
				{Slot: 1, SubSlot: 1, InputSize: 4, OutputSize: 0},
			},
		},
	}

	// 3. Create a new driver instance.
	// Thanks to build tags, this will create a real driver on Linux
	// and a stub driver on other platforms (macOS, Windows).
	driver := profinet.NewDriver(connConfig, deviceConfigs)

	// 4. Start the driver.
	err := driver.Start()
	if err != nil {
		// On non-Linux systems, we expect a specific error.
		if errors.Is(err, profinet.ErrUnsupportedOnPlatform) {
			log.Println("PROFINET is not supported on this platform. This is expected.")
			log.Println("The stub adapter worked correctly. Example finished.")
			return // Graceful exit on non-Linux platforms.
		}
		// For any other error, log it and exit.
		log.Fatalf("Failed to start PROFINET driver: %v", err)
	}

	// The following code will only run on Linux where the driver starts successfully.
	log.Println("PROFINET driver started successfully on Linux.")
	defer func() {
		if err := driver.Stop(); err != nil {
			log.Printf("Error stopping driver: %v", err)
		}
	}()

	// 5. Interact with the driver (e.g., read from cache periodically).
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			cache := driver.GetCache()
			samples := cache.GetAll()
			if len(samples) > 0 {
				log.Printf("--- Reading from Cache (%d samples) ---", len(samples))
				for _, s := range samples {
					log.Printf("  PointID: %s, Value: %v, Quality: %s", s.PointID, s.Value, s.Quality)
				}
			} else {
				log.Println("--- Cache is currently empty ---")
			}
		}
	}()

	// 6. Wait for a shutdown signal to gracefully stop the driver.
	log.Println("Example is running. Press Ctrl+C to stop.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown signal received. Stopping driver...")
}