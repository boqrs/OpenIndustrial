package hj212

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
)

// tcpAdapter implements the Adapter interface for TCP communication.
type tcpAdapter struct {
	mu     sync.Mutex
	config ConnectionConfig
	conn   net.Conn
}

// NewTCPAdapter creates a new adapter for TCP communication.
func NewTCPAdapter() Adapter {
	return &tcpAdapter{}
}

// Connect establishes a TCP connection.
func (a *tcpAdapter) Connect(ctx context.Context, cfg ConnectionConfig) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn != nil {
		return ErrAlreadyConnected
	}

	dialer := net.Dialer{Timeout: cfg.Timeout}
	conn, err := dialer.DialContext(ctx, "tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	a.conn = conn
	a.config = cfg
	return nil
}

// Disconnect closes the TCP connection.
func (a *tcpAdapter) Disconnect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		return nil
	}

	err := a.conn.Close()
	a.conn = nil
	return err
}

// IsConnected returns true if the TCP connection is active.
func (a *tcpAdapter) IsConnected() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.conn != nil
}

// ReadDataSegment listens for incoming data, finds a complete HJ/T 212 frame,
// decodes it, and parses it into a DataSegment.
func (a *tcpAdapter) ReadDataSegment(ctx context.Context) (*DataSegment, error) {
	if !a.IsConnected() {
		return nil, ErrNotConnected
	}

	// Use a buffered reader for efficient scanning.
	reader := bufio.NewReader(a.conn)

	// Scan until we find the start of a message "##".
	_, err := reader.ReadBytes('#')
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadFailed, err)
	}
	second, err := reader.ReadByte()
	if err != nil || second != '#' {
		return nil, fmt.Errorf("%w: could not find start of frame", ErrReadFailed)
	}

	// We have the header. Now read the rest of the potential message.
	// This is a simplified approach. A more robust solution would handle streaming data better.
	buffer := make([]byte, 4096) // Max message size
	buffer[0], buffer[1] = '#', '#'
	n, err := reader.Read(buffer[2:])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadFailed, err)
	}

	// Decode the raw frame to get the CP data.
	cpDataBytes, err := DecodeFrame(buffer[:n+2])
	if err != nil {
		return nil, err
	}

	// Parse the CP data into a structured segment.
	segment, err := ParseDataSegment(string(cpDataBytes))
	if err != nil {
		return nil, err
	}

	return segment, nil
}

// SendCommand builds a complete HJ/T 212 frame from a DataSegment and sends it over TCP.
func (a *tcpAdapter) SendCommand(ctx context.Context, segment *DataSegment) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.IsConnected() {
		return ErrNotConnected
	}

	// Build the CP data string.
	cpData := BuildDataSegment(segment)

	// Encode the CP data into a full HJ/T 212 frame.
	frameBytes, err := EncodeFrame([]byte(cpData))
	if err != nil {
		return err
	}

	// Send the frame.
	_, err = a.conn.Write(frameBytes)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	return nil
}