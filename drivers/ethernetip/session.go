package ethernetip

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// encapsulationHeader is the header for all EtherNet/IP messages.
type encapsulationHeader struct {
	Command       uint16
	Length        uint16
	SessionHandle uint32
	Status        uint32
	SenderContext [8]byte
	Options       uint32
}

const (
	headerSize = 24
)

// Command codes for the encapsulation layer.
const (
	cmdRegisterSession   = 0x0065
	cmdUnRegisterSession = 0x0066
	cmdSendRRData        = 0x006F // For Unconnected Messages (used in Explicit Messaging)
	cmdSendUnitData      = 0x0070 // For Connected Messages (used in Implicit Messaging)
)

// Session manages the TCP connection and session handle with an EtherNet/IP device.
type Session struct {
	conn   net.Conn
	Handle uint32
	config ConnectionConfig
}

// NewSession creates a new session object.
func NewSession(cfg ConnectionConfig) *Session {
	return &Session{
		config: cfg,
	}
}

// Connect dials the device and registers a new session.
func (s *Session) Connect() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	conn, err := net.DialTimeout("tcp", addr, s.config.Timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}
	s.conn = conn
	return s.register()
}

// register sends a RegisterSession command and stores the returned handle.
func (s *Session) register() error {
	header := encapsulationHeader{
		Command: cmdRegisterSession,
		Length:  4, // Protocol version (2 bytes) + options (2 bytes)
	}
	// The payload for RegisterSession is: Protocol Version (uint16) and Options (uint16)
	payload := []byte{1, 0, 0, 0} // Protocol version 1, options 0

	if err := s.send(header, payload); err != nil {
		return fmt.Errorf("failed to send register session request: %w", err)
	}

	respHeader, _, err := s.receive()
	if err != nil {
		return fmt.Errorf("failed to receive register session response: %w", err)
	}

	if respHeader.Status != 0 {
		return fmt.Errorf("register session failed with EIP status code: 0x%X", respHeader.Status)
	}

	s.Handle = respHeader.SessionHandle
	return nil
}

// Close unregisters the session and closes the connection.
func (s *Session) Close() error {
	if s.conn == nil {
		return nil
	}
	defer s.conn.Close()

	header := encapsulationHeader{
		Command:       cmdUnRegisterSession,
		SessionHandle: s.Handle,
	}
	// No payload is needed for unregistering.
	_ = s.send(header, nil) // We attempt to unregister but don't fail if it doesn't work.
	return nil
}

// send writes a command and payload to the connection.
func (s *Session) send(header encapsulationHeader, payload []byte) error {
	if s.conn == nil {
		return ErrNotConnected
	}
	header.Length = uint16(len(payload))
	if err := binary.Write(s.conn, binary.LittleEndian, header); err != nil {
		return err
	}
	if len(payload) > 0 {
		if _, err := s.conn.Write(payload); err != nil {
			return err
		}
	}
	return nil
}

// receive reads a command and payload from the connection.
func (s *Session) receive() (encapsulationHeader, []byte, error) {
	var header encapsulationHeader
	if err := binary.Read(s.conn, binary.LittleEndian, &header); err != nil {
		return header, nil, err
	}

	if header.Length > 0 {
		payload := make([]byte, header.Length)
		if _, err := io.ReadFull(s.conn, payload); err != nil {
			return header, nil, err
		}
		return header, payload, nil
	}

	return header, nil, nil
}