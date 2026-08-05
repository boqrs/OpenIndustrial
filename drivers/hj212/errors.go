package hj212

import "errors"

// Standard HJ/T 212 driver errors
var (
	ErrInvalidHeader         = errors.New("invalid message header (must be ##)")
	ErrInvalidFooter         = errors.New("invalid message footer (must be &&)")
	ErrInvalidLength         = errors.New("message length field does not match actual length")
	ErrCRCMismatch           = errors.New("crc checksum mismatch")
	ErrDataSegmentParse      = errors.New("failed to parse data segment (CP field)")
	ErrMissingRequiredField  = errors.New("a required field (e.g., QN, MN) is missing")
	ErrUnsupportedCommand    = errors.New("unsupported command (CN)")
	ErrConnectionFailed      = errors.New("connection failed")
	ErrNotConnected          = errors.New("not connected")
	ErrAlreadyConnected      = errors.New("already connected")
	ErrSendFailed            = errors.New("send operation failed")
	ErrReadFailed            = errors.New("read operation failed")
)