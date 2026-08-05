package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/drivers/ethercat"
)

func main() {
	log.Println("Starting EtherCAT driver example...")

	// 1. Define the connection configuration for the EtherCAT master.
	// IMPORTANT: "eth0" must be replaced with the actual dedicated network interface
	// for the EtherCAT bus. This interface cannot be used for other TCP/IP traffic.
	connConfig := &ethercat.ConnectionConfig{
		Interface: "eth0",
		CycleTime: 2 * time.Millisecond, // 2ms cycle time is common
	}

	// 2. Define the Process Data Object (PDO) mapping.
	// This is the most critical part of the configuration. You need to know
	// which data (Index, SubIndex) from which slave you want to read.
	// This information comes from the slave's documentation (e.g., ESI file).
	pdoMap := []ethercat.PDOMapping{
		// Example: Reading the status word and actual position from a servo drive at slave index 1.
		{
			Slave:     1,
			Direction: ethercat.PDOInput, // We are reading data from the slave.
			Entries: []ethercat.PDOEntry{
				{
					Name:      "servo1_status_word",
					Index:     0x6041, // CiA 402 Statusword
					SubIndex:  0,
					Offset:    0,
					BitLength: 16,
					DataType:  ethercat.TypeUINT,
				},
				{
					Name:      "servo1_actual_position",
					Index:     0x6064, // CiA 402 Position actual value
					SubIndex:  0,
					Offset:    16, // Starts after the 16-bit status word
					BitLength: 32,
					DataType:  ethercat.TypeDINT,
				},
			},
		},
		// Example: Reading 8 digital inputs from an I/O terminal at slave index 2.
		{
			Slave:     2,
			Direction: ethercat.PDOInput,
			Entries: []ethercat.PDOEntry{
				{Name: "di_1", Index: 0x6000, SubIndex: 1, Offset: 0, BitLength: 1, DataType: ethercat.TypeBOOL},
				{Name: "di_2", Index: 0x6000, SubIndex: 2, Offset: 1, BitLength: 1, DataType: ethercat.TypeBOOL},
				{Name: "di_3", Index: 0x6000, SubIndex: 3, Offset: 2, BitLength: 1, DataType: ethercat.TypeBOOL},
				{Name: "di_4", Index: 0x6000, SubIndex: 4, Offset: 3, BitLength: 1, DataType: ethercat.TypeBOOL},
			},
		},
	}

	// 3. Create a new driver instance.
	driver := ethercat.NewDriver(connConfig, pdoMap)

	// 4. Start the driver.
	// This will initialize SOEM, scan the bus, configure slaves, and start the real-time cycle.
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start EtherCAT driver: %v", err)
	}
	defer driver.Stop()

	// 5. Start a goroutine to consume samples from the driver.
	go func() {
		log.Println("Starting sample consumer...")
		sampleChan := driver.Samples()
		for sample := range sampleChan {
			log.Printf("--> Received Sample: PointID=[%s], Value=[%v], Quality=[%s]",
				sample.PointID,
				sample.Value,
				sample.Quality,
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