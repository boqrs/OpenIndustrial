package dlt645

import "errors"

var (
	ErrInvalidFrameLength      = errors.New("invalid frame length")
	ErrInvalidFrameStartOrEnd  = errors.New("invalid frame start or end bytes")
	ErrChecksumMismatch        = errors.New("frame checksum mismatch")
	ErrInvalidAddressFormat    = errors.New("invalid meter address format, must be 12 BCD characters")
	ErrUnsupportedDataType     = errors.New("unsupported data type for decoding")
	ErrDataLengthMismatch      = errors.New("data length does not match expected length for the data type")
	ErrPollingNotStarted       = errors.New("poller has not been started")
	ErrConnectionFailed        = errors.New("connection failed")
	ErrNotConnected            = errors.New("not connected")
	ErrAlreadyConnected        = errors.New("already connected")
	ErrWriteFailed             = errors.New("write operation failed")
	ErrReadFailed              = errors.New("read operation failed")
)