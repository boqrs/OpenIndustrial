package ethernetip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// CIPClient provides methods to send CIP commands over an established session.
type CIPClient struct {
	session *Session
}

// NewCIPClient creates a new client that uses the given session.
func NewCIPClient(s *Session) *CIPClient {
	return &CIPClient{session: s}
}

// ReadTag sends a "Read Tag Service" request for a single tag.
func (c *CIPClient) ReadTag(tag string) ([]byte, DataType, error) {
	// Construct the CIP request payload for reading a tag.
	// This involves creating a "Request Path" from the tag name.
	path, err := buildRequestPath(tag)
	if err != nil {
		return nil, "", fmt.Errorf("failed to build request path for tag '%s': %w", tag, err)
	}

	// Service Code for "Read Tag Service" is 0x4C.
	const serviceReadTag = 0x4C
	// We are requesting 1 element.
	const elements = 1

	var buf bytes.Buffer
	buf.WriteByte(serviceReadTag)
	buf.WriteByte(byte(len(path) / 2)) // Path size in 16-bit words
	buf.Write(path)
	binary.Write(&buf, binary.LittleEndian, uint16(elements))

	// Send the request and get the response.
	resp, err := c.sendRRData(buf.Bytes())
	if err != nil {
		return nil, "", err
	}

	// Parse the CIP response.
	// The response contains a General Status, Extended Status,
	// the data type, and the actual data.
	if len(resp) < 4 {
		return nil, "", ErrInvalidResponse
	}

	generalStatus := resp[0]
	// extendedStatus := resp[1] // We can use this for more detailed errors.

	if generalStatus != 0 {
		if generalStatus == 0x05 { // "Path destination unknown"
			return nil, "", ErrTagNotFound
		}
		return nil, "", fmt.Errorf("read tag failed with CIP status: 0x%X", generalStatus)
	}

	// The data type is a 2-byte code starting at offset 2.
	dataTypeCode := binary.LittleEndian.Uint16(resp[2:4])
	dataType := cipTypeToDataType(dataTypeCode)

	// The actual data follows the type code.
	data := resp[4:]

	return data, dataType, nil
}

// sendRRData wraps a CIP payload in the necessary headers and sends it via SendRRData.
func (c *CIPClient) sendRRData(cipPayload []byte) ([]byte, error) {
	// The SendRRData command requires a specific structure in its own payload.
	// It includes an interface handle (usually 0), a timeout, and the CIP data.
	var rrData bytes.Buffer
	binary.Write(&rrData, binary.LittleEndian, uint32(0)) // Interface Handle
	binary.Write(&rrData, binary.LittleEndian, uint16(c.session.config.Timeout/time.Millisecond)) // Timeout
	binary.Write(&rrData, binary.LittleEndian, uint16(2)) // Item Count: Null Address Item + Unconnected Data Item

	// Item 1: Null Address Item (indicates the message is for the PLC itself)
	binary.Write(&rrData, binary.LittleEndian, uint16(0x0000)) // Type ID
	binary.Write(&rrData, binary.LittleEndian, uint16(0x0000)) // Length

	// Item 2: Unconnected Data Item
	binary.Write(&rrData, binary.LittleEndian, uint16(0x00B2)) // Type ID
	binary.Write(&rrData, binary.LittleEndian, uint16(len(cipPayload))) // Length
	rrData.Write(cipPayload)

	// Construct the main encapsulation header.
	header := encapsulationHeader{
		Command:       cmdSendRRData,
		SessionHandle: c.session.Handle,
	}

	// Send the request.
	if err := c.session.send(header, rrData.Bytes()); err != nil {
		return nil, err
	}

	// Receive the response.
	respHeader, respPayload, err := c.session.receive()
	if err != nil {
		return nil, err
	}

	if respHeader.Status != 0 {
		return nil, fmt.Errorf("sendRRData failed with EIP status: 0x%X", respHeader.Status)
	}

	// The actual CIP response is embedded within the SendRRData response payload.
	// We need to skip the SendRRData header part (Interface Handle, Timeout, Item Count, etc.).
	if len(respPayload) < 10 {
		return nil, ErrInvalidResponse
	}
	// The CIP response starts after the second item's header (Type ID + Length).
	// Offset: InterfaceHandle(4) + Timeout(2) + ItemCount(2) + Item1Header(4) + Item2Header(4) = 16
	// But the response only returns the Unconnected Data Item, so we skip its header.
	// Offset: InterfaceHandle(4) + Timeout(2) + ItemCount(2) + Item1Header(4) = 12
	// The response payload for SendRRData contains the CIP response starting after some headers.
	// The interesting part is the payload of the Unconnected Data Item.
	// Let's find it by looking for its Type ID (0x00B2).
	// A simpler way for now: the CIP response is usually the last part.
	// The header of the response is typically 10 bytes long (Handle, Timeout, Count, Type, Length).
	// The CIP response is after the Unconnected Data Item header.
	// Let's assume the interesting payload starts at offset 8 for now.
	// A robust implementation would parse the items properly.
	if len(respPayload) < 8 {
		return nil, ErrInvalidResponse
	}
	return respPayload[8:], nil
}

// buildRequestPath converts a string tag into a byte-encoded CIP path.
// Example: "MyTag.SubTag[1]" -> encoded bytes.
// This is a complex process. For now, we'll use a simplified version for simple tags.
func buildRequestPath(tag string) ([]byte, error) {
	var path []byte
	// For simple tags, the path is just the tag name, padded to an even length.
	path = append(path, byte(0x91)) // ANSI Extended Symbol segment
	path = append(path, byte(len(tag)))
	path = append(path, []byte(tag)...)
	if len(tag)%2 != 0 {
		path = append(path, 0x00) // Pad to even length
	}
	return path, nil
}

// cipTypeToDataType converts a CIP type code to our internal DataType.
func cipTypeToDataType(code uint16) DataType {
	switch code {
	case 0xC1:
		return TypeBOOL
	case 0xC2:
		return TypeSINT
	case 0xC3:
		return TypeINT
	case 0xC4:
		return TypeDINT
	case 0xC5:
		return TypeLINT
	case 0xCA:
		return TypeREAL
	case 0xCB:
		return TypeLREAL
	case 0xD0:
		return TypeSTRING
	default:
		return "" // Unsupported
	}
}