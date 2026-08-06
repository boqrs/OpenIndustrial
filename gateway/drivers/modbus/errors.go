package modbus

import "errors"

var (

	// Driver

	ErrNotInitialized = errors.New("modbus driver not initialized")

	ErrAlreadyStarted = errors.New("modbus driver already started")

	ErrNotStarted = errors.New("modbus driver not started")

	// Connection

	ErrConnectionClosed = errors.New("connection closed")

	ErrConnectionTimeout = errors.New("connection timeout")

	ErrDisconnected = errors.New("device disconnected")

	// Transport

	ErrUnsupportedMode = errors.New("unsupported transport mode")

	ErrTransportUnavailable = errors.New("transport unavailable")

	// Register

	ErrInvalidRegister = errors.New("invalid register")

	ErrInvalidAddress = errors.New("invalid register address")

	ErrInvalidLength = errors.New("invalid register length")

	// Codec

	ErrCodec = errors.New("codec error")

	ErrUnsupportedDataType = errors.New("unsupported data type")

	ErrInvalidEndian = errors.New("invalid endian")

	ErrDecode = errors.New("decode failed")

	ErrEncode = errors.New("encode failed")

	// Mapping

	ErrPointNotFound = errors.New("point mapping not found")

	ErrDuplicatePoint = errors.New("duplicate point")

	ErrDuplicateAddress = errors.New("duplicate register address")

	// Collector

	ErrCollectorStopped = errors.New("collector stopped")

	ErrEmptyBatch = errors.New("empty batch")

	ErrBatchOverflow = errors.New("batch overflow")
)