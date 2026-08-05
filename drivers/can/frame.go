package can

import "time"

// Frame represents a single CAN or CAN FD frame.
type Frame struct {
	// ID is the CAN identifier of the frame.
	ID uint32
	// Extended is true if the frame uses a 29-bit extended identifier.
	Extended bool
	// RTR is true if the frame is a Remote Transmission Request.
	RTR bool
	// DLC is the Data Length Code, indicating the length of the data in bytes.
	DLC uint8
	// Data is the payload of the frame, up to 8 bytes for classic CAN
	// and up to 64 bytes for CAN FD.
	Data []byte
	// Timestamp is the time at which the frame was received.
	Timestamp time.Time
}