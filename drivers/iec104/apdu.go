package iec104

import (
	"encoding/binary"
	"fmt"
)

const (
	APCI_START_BYTE = 0x68
	APCI_MIN_LENGTH = 6
)

// FrameType identifies the type of the APCI frame.
type FrameType int

const (
	FrameTypeI FrameType = iota // Information frame, carries an ASDU
	FrameTypeS                  // Supervisory frame, for acknowledgements
	FrameTypeU                  // Unnumbered frame, for control (STARTDT, STOPDT, TESTFR)
)

// UFrameFunction defines the function of a U-frame.
type UFrameFunction byte

const (
	UFrameStartDTAct UFrameFunction = 0x07
	UFrameStartDTCon UFrameFunction = 0x0B
	UFrameStopDTAct  UFrameFunction = 0x13
	UFrameStopDTCon  UFrameFunction = 0x23
	UFrameTestFRAct  UFrameFunction = 0x43
	UFrameTestFRCon  UFrameFunction = 0x83
)

// APCI (Application Protocol Control Information) represents the header of an IEC 104 frame.
type APCI struct {
	Start      byte
	Length     uint8
	Control    [4]byte
	FrameType  FrameType
	SendSeqNum uint16
	RecvSeqNum uint16
}

// APDU (Application Protocol Data Unit) represents a full IEC 104 message.
type APDU struct {
	APCI APCI
	ASDU []byte // Raw ASDU data
}

// ParseAPDU parses a byte slice into an APDU structure.
func ParseAPDU(data []byte) (*APDU, error) {
	if len(data) < APCI_MIN_LENGTH {
		return nil, fmt.Errorf("invalid APDU length: got %d, want at least %d", len(data), APCI_MIN_LENGTH)
	}
	if data[0] != APCI_START_BYTE {
		return nil, fmt.Errorf("invalid start byte: got 0x%X, want 0x%X", data[0], APCI_START_BYTE)
	}

	length := data[1]
	if len(data) != int(length)+2 {
		return nil, fmt.Errorf("APDU length mismatch: header says %d, actual is %d", length, len(data)-2)
	}

	apci := APCI{
		Start:  data[0],
		Length: length,
		Control: [4]byte{
			data[2], data[3], data[4], data[5],
		},
	}

	// Determine frame type from control field
	if (apci.Control[0] & 0x01) == 0 {
		apci.FrameType = FrameTypeI
		apci.SendSeqNum = binary.LittleEndian.Uint16(apci.Control[0:2]) >> 1
		apci.RecvSeqNum = binary.LittleEndian.Uint16(apci.Control[2:4]) >> 1
	} else if (apci.Control[0] & 0x03) == 1 {
		apci.FrameType = FrameTypeS
		apci.RecvSeqNum = binary.LittleEndian.Uint16(apci.Control[2:4]) >> 1
	} else if (apci.Control[0] & 0x03) == 3 {
		apci.FrameType = FrameTypeU
	}

	asduData := data[APCI_MIN_LENGTH:]

	return &APDU{
		APCI: apci,
		ASDU: asduData,
	}, nil
}

// NewUFrame creates a new U-frame for control purposes.
func NewUFrame(uFunc UFrameFunction) []byte {
	frame := make([]byte, APCI_MIN_LENGTH)
	frame[0] = APCI_START_BYTE
	frame[1] = 4 // Length of control field
	frame[2] = byte(uFunc)
	frame[3] = 0
	frame[4] = 0
	frame[5] = 0
	return frame
}

// NewIFrame creates a new I-frame to carry an ASDU.
func NewIFrame(sendSeq, recvSeq uint16, asdu []byte) []byte {
	apduLength := 4 + len(asdu) // 4 bytes for control field + ASDU length
	frame := make([]byte, apduLength+2) // +2 for start and length bytes

	frame[0] = APCI_START_BYTE
	frame[1] = byte(apduLength)

	// Encode sequence numbers
	binary.LittleEndian.PutUint16(frame[2:4], sendSeq<<1)
	binary.LittleEndian.PutUint16(frame[4:6], recvSeq<<1)

	// Copy ASDU
	copy(frame[6:], asdu)

	return frame
}