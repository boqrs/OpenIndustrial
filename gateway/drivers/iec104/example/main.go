package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway/drivers/iec104"
)

// ExampleDriver demonstrates the full lifecycle of the IEC 104 driver,
// including configuration, starting, reading, writing, and stopping.
func main() {
	// Start a mock IEC 104 server to simulate a device.
	serverAddr := "127.0.0.1:2404"
	go startMockServer(serverAddr)
	time.Sleep(500 * time.Millisecond) // Give the server a moment to start.

	// 1. Configuration
	config := &iec104.ConnectionConfig{
		Host:          "127.0.0.1",
		Port:          2404,
		Timeout:       5 * time.Second,
		CommonAddress: 1,
	}

	points := []iec104.PointMapping{
		{ID: "TI-101", IOA: 4100, DataType: "M_ME_NC_1"},
		{ID: "QI-102", IOA: 10, DataType: "M_SP_NA_1"},
		{ID: "CMD-01", IOA: 6100, DataType: "C_SC_NA_1"},
	}

	// 2. Initialization
	driver, err := iec104.NewDriver(config, points)
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	// Start the driver.
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver: %v", err)
	}
	defer driver.Stop()

	// Set up a channel for graceful shutdown.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("Driver started. Will run for 15 seconds. Press Ctrl+C to exit earlier.")

	// 3. Main application loop
loop:
	for {
		select {
		case <-time.After(5 * time.Second):
			// 3a. Data Reading
			sample, err := driver.Read("TI-101")
			if err != nil {
				fmt.Printf("Failed to read point 'TI-101': %v\n", err)
			} else {
				fmt.Printf("Read data: PointID=%s, Value=%.2f\n", sample.PointID, sample.Value)
			}

			// 3b. Data Writing
			fmt.Println("Sending a single command 'true' to CMD-01...")
			err = driver.Write("CMD-01", true)
			if err != nil {
				fmt.Printf("Failed to write to point 'CMD-01': %v\n", err)
			}

		case <-stopChan:
			log.Println("Shutdown signal received.")
			break loop
		case <-ctx.Done():
			log.Println("Example finished after 15 seconds.")
			break loop
		}
	}

	// Unordered output:
	// Driver started. Will run for 15 seconds. Press Ctrl+C to exit earlier.
	// Read data: PointID=TI-101, Value=98.60
	// Sending a single command 'true' to CMD-01...
	// Read data: PointID=TI-101, Value=98.60
	// Sending a single command 'true' to CMD-01...
	// Example finished after 15 seconds.
}

// startMockServer simulates a basic IEC 104 slave device.
func startMockServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// If the port is already in use, we assume another test is running the server.
		if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "listen: address already in use" {
			log.Printf("Mock server address %s is already in use, assuming it's running.", address)
			return
		}
		log.Fatalf("Mock server failed to listen: %v", err)
	}
	defer listener.Close()
	log.Printf("Mock server listening on %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return // Listener closed
		}
		go handleMockConnection(conn)
	}
}

// handleMockConnection handles a single client connection for the mock server.
func handleMockConnection(conn net.Conn) {
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Read STARTDT_ACT
	buf := make([]byte, 256)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	apdu, err := iec104.ParseAPDU(buf[:n])
	if err != nil || apdu.APCI.FrameType != iec104.FrameTypeU || apdu.APCI.Control[0] != byte(iec104.UFrameStartDTAct) {
		return
	}

	// Respond with STARTDT_CON
	startDTConFrame := iec104.NewUFrame(iec104.UFrameStartDTCon)
	conn.Write(startDTConFrame)

	// Wait for General Interrogation
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	conn.Read(buf)

	// Respond with activation confirmation for interrogation
	asduInterrogationConf, _ := iec104.MarshalASDU(&iec104.ASDU{
		TypeID:        iec104.C_IC_NA_1,
		VSQ:           1,
		Cause:         iec104.COT_ACTIVATION_CON,
		CommonAddress: 1,
		Objects:       []iec104.InformationObject{{IOA: 0, Value: 20}},
	})
	iFrameConf := iec104.NewIFrame(0, 1, asduInterrogationConf)
	conn.Write(iFrameConf)

	// Send some initial data
	sendData(conn, 0, 2)

	// End of interrogation
	asduInterrogationEnd, _ := iec104.MarshalASDU(&iec104.ASDU{
		TypeID:        iec104.C_IC_NA_1,
		VSQ:           1,
		Cause:         iec104.COT_ACTIVATION_TERMINATION,
		CommonAddress: 1,
		Objects:       []iec104.InformationObject{{IOA: 0, Value: 20}},
	})
	iFrameEnd := iec104.NewIFrame(1, 2, asduInterrogationEnd)
	conn.Write(iFrameEnd)

	// Periodically send spontaneous data
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	var sendSeq uint16 = 2
	var recvSeq uint16 = 2

	for {
		select {
		case <-ticker.C:
			sendData(conn, sendSeq, recvSeq)
			sendSeq++
		case <-ctx.Done():
			return
		}
	}
}

// sendData constructs and sends a sample data frame.
func sendData(conn net.Conn, sendSeq, recvSeq uint16) {
	asduData, _ := iec104.MarshalASDU(&iec104.ASDU{
		TypeID:        iec104.M_ME_NC_1,
		VSQ:           1,
		Cause:         iec104.COT_SPONTANEOUS,
		CommonAddress: 1,
		Objects:       []iec104.InformationObject{{IOA: 4100, Value: 98.6, Timestamp: time.Now()}},
	})
	iFrameData := iec104.NewIFrame(sendSeq, recvSeq, asduData)
	conn.Write(iFrameData)
}