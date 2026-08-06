package iec104

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// tcpAdapter implements the Adapter interface using a standard TCP connection.
type tcpAdapter struct {
	mu         sync.RWMutex
	conn       net.Conn
	config     ConnectionConfig
	connected  bool
	cancelRead context.CancelFunc
	wg         sync.WaitGroup

	// Sequence numbers
	sendSeqNum uint16
	recvSeqNum uint16
}

// NewTCPAdapter creates a new adapter for TCP-based IEC 104 communication.
func NewTCPAdapter() Adapter {
	return &tcpAdapter{}
}

func (a *tcpAdapter) Connect(ctx context.Context, cfg ConnectionConfig) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.connected {
		return ErrAlreadyConnected
	}

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	dialer := net.Dialer{Timeout: cfg.Timeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	a.conn = conn
	a.config = cfg

	// Send STARTDT and wait for confirmation
	if err := a.sendUFrame(UFrameStartDTAct); err != nil {
		a.conn.Close()
		return fmt.Errorf("failed to send STARTDT: %w", err)
	}

	// Simple wait for STARTDT_CON, a more robust implementation would have a timeout
	// and proper frame parsing here.
	buf := make([]byte, 255)
	a.conn.SetReadDeadline(time.Now().Add(a.config.Timeout))
	n, err := conn.Read(buf)
	if err != nil {
		a.conn.Close()
		return fmt.Errorf("failed to read STARTDT confirmation: %w", err)
	}
	apdu, err := ParseAPDU(buf[:n])
	if err != nil || apdu.APCI.FrameType != FrameTypeU || apdu.APCI.Control[0] != byte(UFrameStartDTCon) {
		a.conn.Close()
		return fmt.Errorf("did not receive STARTDT confirmation")
	}

	a.connected = true
	return nil
}

func (a *tcpAdapter) Disconnect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.connected {
		return nil
	}

	// Stop the reading loop
	if a.cancelRead != nil {
		a.cancelRead()
	}
	a.wg.Wait() // Wait for the loop to finish

	// Send STOPDT
	_ = a.sendUFrame(UFrameStopDTAct) // Best effort

	err := a.conn.Close()
	a.connected = false
	a.conn = nil
	return err
}

func (a *tcpAdapter) IsConnected() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.connected
}

func (a *tcpAdapter) Read(ctx context.Context, points []PointMapping) ([]Sample, error) {
	return nil, fmt.Errorf("targeted Read not implemented, use Subscribe")
}

func (a *tcpAdapter) Write(ctx context.Context, point PointMapping, value interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.connected {
		return ErrNotConnected
	}

	// For simplicity, we'll map a boolean write to a Single Command (C_SC_NA_1)
	// A more robust implementation would use the PointMapping to determine the correct TypeID.
	val, ok := value.(bool)
	if !ok {
		return fmt.Errorf("currently only boolean writes are supported")
	}

	asdu := &ASDU{
		TypeID:        C_SC_NA_1,
		VSQ:           1, // 1 object
		Cause:         COT_ACTIVATION,
		CommonAddress: a.config.CommonAddress,
		Objects: []InformationObject{
			{
				IOA:   point.IOA,
				Value: val,
			},
		},
	}

	asduBytes, err := MarshalASDU(asdu)
	if err != nil {
		return fmt.Errorf("failed to marshal write command: %w", err)
	}

	// Wrap in I-frame and send
	iFrame := NewIFrame(a.sendSeqNum, a.recvSeqNum, asduBytes)
	a.conn.SetWriteDeadline(time.Now().Add(a.config.Timeout))
	if _, err := a.conn.Write(iFrame); err != nil {
		return fmt.Errorf("failed to write command to connection: %w", err)
	}

	a.sendSeqNum++
	return nil
}

func (a *tcpAdapter) Subscribe(ctx context.Context, ch chan<- InformationObject) error {
	if !a.IsConnected() {
		return ErrNotConnected
	}

	readCtx, cancel := context.WithCancel(ctx)
	a.cancelRead = cancel
	a.wg.Add(1)

	// Start the background reading loop
	go a.readLoop(readCtx, ch)

	// Send a General Interrogation command
	if err := a.handleInterrogation(ctx); err != nil {
		// We might still be able to receive spontaneous data, so don't necessarily fail here.
		fmt.Printf("Warning: failed to send interrogation command: %v\n", err)
	}

	// The readLoop will now handle incoming data.
	// The function will block until the context is cancelled.
	<-readCtx.Done()
	return readCtx.Err()
}

// sendUFrame sends a U-frame with the given control command.
func (a *tcpAdapter) sendUFrame(command UFrameFunction) error {
	frame := NewUFrame(command)
	_, err := a.conn.Write(frame)
	return err
}

// handleInterrogation sends a general interrogation command.
func (a *tcpAdapter) handleInterrogation(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	asdu := &ASDU{
		TypeID:        C_IC_NA_1,
		VSQ:           1, // 1 object
		Cause:         COT_ACTIVATION,
		CommonAddress: a.config.CommonAddress,
		Objects: []InformationObject{
			{
				IOA:   0,    // IOA 0 means station interrogation
				Value: 20, // QOI = 20 means "interrogate station"
			},
		},
	}

	asduBytes, err := MarshalASDU(asdu)
	if err != nil {
		return fmt.Errorf("failed to marshal interrogation command: %w", err)
	}

	// Wrap in I-frame and send
	iFrame := NewIFrame(a.sendSeqNum, a.recvSeqNum, asduBytes)
	a.conn.SetWriteDeadline(time.Now().Add(a.config.Timeout))
	if _, err := a.conn.Write(iFrame); err != nil {
		return fmt.Errorf("failed to write interrogation command to connection: %w", err)
	}

	a.sendSeqNum++
	return nil
}

// readLoop continuously reads from the connection and processes frames.
func (a *tcpAdapter) readLoop(ctx context.Context, ch chan<- InformationObject) {
	defer a.wg.Done()
	buf := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			a.conn.SetReadDeadline(time.Now().Add(a.config.Timeout))
			n, err := a.conn.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Read error: %v\n", err)
				}
				return // End loop on read error or EOF
			}

			if n > 0 {
				a.processFrame(buf[:n], ch)
			}
		}
	}
}

// processFrame parses the incoming byte slice and handles the APDU frame.
func (a *tcpAdapter) processFrame(data []byte, ch chan<- InformationObject) {
	apdu, err := ParseAPDU(data)
	if err != nil {
		fmt.Printf("APDU parsing error: %v\n", err)
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	switch apdu.APCI.FrameType {
	case FrameTypeI:
		// Check sequence numbers
		if apdu.APCI.RecvSeqNum != a.sendSeqNum {
			// Handle sequence error
		}
		if apdu.APCI.SendSeqNum != a.recvSeqNum {
			// Handle sequence error
		}
		a.recvSeqNum++

		asdu, err := ParseASDU(apdu.ASDU)
		if err != nil {
			fmt.Printf("ASDU parsing error: %v\n", err)
			return
		}

		// Send InformationObjects to channel
		for _, obj := range asdu.Objects {
			ch <- obj
		}

	case FrameTypeS:
		// Handle S-Frame (acknowledgement)
		if apdu.APCI.RecvSeqNum != a.sendSeqNum {
			// Handle sequence error
		}
	case FrameTypeU:
		// Handle U-Frame (e.g., TESTFR_CON)
	}
}