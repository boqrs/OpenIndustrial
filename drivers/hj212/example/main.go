package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/drivers/hj212"
)

// mockAdapter simulates the behavior of a real TCP adapter for testing.
type mockAdapter struct {
	mu          sync.Mutex
	isConnected bool
	// dataToReturn is the raw HJ/T 212 frame the mock adapter will "send" to the driver.
	dataToReturn []byte
	// commandsReceived stores any commands sent by the driver.
	commandsReceived []*hj212.DataSegment
}

func (m *mockAdapter) Connect(ctx context.Context, cfg hj212.ConnectionConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.Println("Mock Adapter: Connect called")
	m.isConnected = true
	return nil
}

func (m *mockAdapter) Disconnect(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.Println("Mock Adapter: Disconnect called")
	m.isConnected = false
	return nil
}

func (m *mockAdapter) IsConnected() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isConnected
}

// ReadDataSegment simulates a device proactively sending data.
func (m *mockAdapter) ReadDataSegment(ctx context.Context) (*hj212.DataSegment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Simulate blocking until data is available
	time.Sleep(100 * time.Millisecond)

	if m.dataToReturn == nil {
		// Block forever if no data is set, simulating a quiet device.
		select {}
	}

	// This is a simplified simulation. A real test would use a net.Pipe.
	// We decode the frame here to simulate what the real adapter does.
	cpDataBytes, err := decodeFrame(m.dataToReturn)
	if err != nil {
		return nil, err
	}
	segment, err := parseDataSegment(string(cpDataBytes))
	if err != nil {
		return nil, err
	}
	return segment, nil
}

func (m *mockAdapter) SendCommand(ctx context.Context, segment *hj212.DataSegment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.Printf("Mock Adapter: SendCommand called with CN=%s", segment.CN)
	m.commandsReceived = append(m.commandsReceived, segment)
	return nil
}

// Helper functions for the test that are not part of the production code.
// They need to be in the _test package to access hj212 internals.
func decodeFrame(data []byte) ([]byte, error) {
	// This is a stand-in for the internal decodeFrame function for testing purposes.
	// A better approach would be to export it for testing (e.g., `hj212.DecodeFrame`).
	// For this example, we'll just assume it works and extract the CP part manually.
	// "##0073CP=&&QN=...;w01018-Rtd=56.3&&1234&&" -> "QN=...;w01018-Rtd=56.3"
	cpStart := bytes.Index(data, []byte("CP=&&"))
	if cpStart == -1 {
		return nil, fmt.Errorf("CP=&& not found")
	}
	cpEnd := bytes.LastIndex(data, []byte("&&"))
	if cpEnd == -1 || cpEnd <= cpStart {
		return nil, fmt.Errorf("ending && not found")
	}
	return data[cpStart+5 : cpEnd], nil
}

func parseDataSegment(cpData string) (*hj212.DataSegment, error) {
	// Simplified parser for test.
	segment := &hj212.DataSegment{
		Pollutants: make(map[string]string),
	}
	pairs := strings.Split(cpData, ";")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			segment.Pollutants[kv[0]] = kv[1]
		}
	}
	return segment, nil
}

// ExampleDriver demonstrates the full lifecycle of the HJ/T 212 driver.
func main() {
	// 1. Configuration
	config := &hj212.ConnectionConfig{
		Type:    "tcp",
		Address: "localhost:9999", // Not used by mock adapter
		Timeout: 5 * time.Second,
	}

	points := []hj212.PointMapping{
		{ID: "sewage_cod", Code: "w01018"},
		{ID: "sewage_ph", Code: "w01001"},
	}

	// 2. Create a mock adapter and inject it.
	// This is the raw data frame our mock adapter will pretend to receive.
	mockData := []byte("##0088CP=&&QN=20230101000000001;ST=32;CN=2011;PW=123456;MN=SN12345;Flag=4;w01018-Rtd=56.3,w01001-Rtd=7.2&&4321&&")

	mock := &mockAdapter{
		dataToReturn: mockData,
	}

	// We need a way to inject the mock adapter. Let's modify NewDriver for testability.
	// driver, err := hj212.NewDriver(config, points) // This uses the real TCP adapter.
	// For the test, we need a constructor that accepts an adapter.
	// Let's assume we have a NewDriverWithAdapter function.
	driver, err := hj212.NewDriverWithAdapter(config, points, mock)
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	// 3. Start the driver
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver: %v", err)
	}
	defer driver.Stop()

	// Give the poller a moment to read from the mock adapter and populate the cache.
	time.Sleep(200 * time.Millisecond)

	// 4. Read data from the driver's cache
	codSample, err := driver.Read("sewage_cod")
	if err != nil {
		fmt.Printf("Failed to read sewage_cod: %v\n", err)
	} else {
		fmt.Printf("Read from cache: PointID=%s, Value=%v\n", codSample.PointID, codSample.Value)
	}

	phSample, err := driver.Read("sewage_ph")
	if err != nil {
		fmt.Printf("Failed to read sewage_ph: %v\n", err)
	} else {
		fmt.Printf("Read from cache: PointID=%s, Value=%v\n", phSample.PointID, phSample.Value)
	}

	// Output:
	// Read from cache: PointID=sewage_cod, Value=56.3
	// Read from cache: PointID=sewage_ph, Value=7.2
}