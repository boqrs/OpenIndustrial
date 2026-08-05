package dlt645

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/goburrow/serial"
)

// serialAdapter implements the Adapter interface for serial (RS485) communication.
type serialAdapter struct {
	mu     sync.Mutex
	config ConnectionConfig
	port   serial.Port
}

// NewSerialAdapter creates a new adapter for serial communication.
func NewSerialAdapter() Adapter {
	return &serialAdapter{}
}

// Connect establishes a serial port connection.
func (a *serialAdapter) Connect(ctx context.Context, cfg ConnectionConfig) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.port != nil {
		return ErrAlreadyConnected
	}

	serialConfig := serial.Config{
		Address:  cfg.Address,
		BaudRate: cfg.BaudRate,
		DataBits: cfg.DataBits,
		StopBits: cfg.StopBits,
		Parity:   cfg.Parity,
		Timeout:  cfg.Timeout,
	}

	port, err := serial.Open(&serialConfig)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	a.port = port
	a.config = cfg
	return nil
}

// Disconnect closes the serial port connection.
func (a *serialAdapter) Disconnect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.port == nil {
		return nil
	}

	err := a.port.Close()
	a.port = nil
	return err
}

// IsConnected returns true if the serial port is open.
func (a *serialAdapter) IsConnected() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.port != nil
}

// Read sends a read request to the serial port and waits for a response.
func (a *serialAdapter) Read(ctx context.Context, meterAddress string, points []PointMapping) ([]Sample, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.IsConnected() {
		return nil, ErrNotConnected
	}

	var samples []Sample

	// For simplicity, we read one point at a time. A real-world optimization
	// might group requests if the device supports it.
	for _, point := range points {
		// 1. Build the request frame
		requestFrame, err := buildReadFrame(meterAddress, point.DI)
		if err != nil {
			return nil, fmt.Errorf("failed to build read frame for DI %d: %w", point.DI, err)
		}

		// 2. Send the frame
		if _, err := a.port.Write(requestFrame); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrWriteFailed, err)
		}

		// 3. Read the response
		// A robust implementation needs a better way to determine the end of a frame
		// than just a fixed-size buffer and a timeout.
		responseBytes := make([]byte, 256)
		n, err := a.port.Read(responseBytes)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrReadFailed, err)
		}
		if n == 0 {
			continue // Or return timeout error
		}

		// 4. Decode the response frame
		responseFrame, err := decodeResponseFrame(responseBytes[:n])
		if err != nil {
			return nil, fmt.Errorf("failed to decode response frame: %w", err)
		}

		// 5. Parse the data payload
		value, err := parseData(responseFrame.Data, point)
		if err != nil {
			return nil, fmt.Errorf("failed to parse data for DI %d: %w", point.DI, err)
		}

		// 6. Create a Sample and add it to the slice
		sample := Sample{
			PointID:   fmt.Sprintf("%s.%s", meterAddress, point.ID),
			Value:     value,
			Timestamp: time.Now(),
			Quality:   QualityGood,
			Source:    "dlt645",
		}
		samples = append(samples, sample)
	}

	return samples, nil
}

// Write sends a write request to the serial port.
// This is a placeholder.
func (a *serialAdapter) Write(ctx context.Context, req WriteRequest) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.IsConnected() {
		return ErrNotConnected
	}

	// In a real implementation, we would:
	// 1. Build a DL/T 645 write frame based on the WriteRequest.
	// 2. Send the frame via a.port.Write().
	// 3. Optionally, wait for and validate the confirmation response.

	return fmt.Errorf("write not yet implemented")
}