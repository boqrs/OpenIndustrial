// drivers/bacnet/errors.go
package bacnet

import "errors"

var (
	// Driver related errors
	ErrNotInitialized = errors.New("bacnet driver not initialized")
	ErrAlreadyStarted = errors.New("bacnet driver already started")
	ErrNotStarted     = errors.New("bacnet driver not started")

	// Connection related errors
	ErrConnectionFailed    = errors.New("bacnet connection failed")
	ErrConnectionClosed    = errors.New("bacnet connection closed")
	ErrConnectionTimeout   = errors.New("bacnet connection timeout")
	ErrDisconnected        = errors.New("bacnet device disconnected")
	ErrUnsupportedMode     = errors.New("unsupported bacnet connection mode")
	ErrDeviceNotFound      = errors.New("bacnet device not found")

	// Object/Property related errors
	ErrObjectNotFound      = errors.New("bacnet object not found")
	ErrPropertyNotFound    = errors.New("bacnet property not found")
	ErrInvalidObject       = errors.New("invalid bacnet object type or instance")
	ErrInvalidProperty     = errors.New("invalid bacnet property ID")
	ErrUnsupportedDataType = errors.New("unsupported bacnet data type for property")
	ErrPropertyNotWritable = errors.New("bacnet property not writable")
	ErrReadFailed          = errors.New("bacnet read failed")
	ErrWriteFailed         = errors.New("bacnet write failed")

	// Configuration related errors
	ErrInvalidConfig       = errors.New("invalid bacnet configuration")
	ErrDuplicateNode       = errors.New("duplicate bacnet node mapping") // e.g., same object/property mapped twice
	ErrDuplicatePointID    = errors.New("duplicate point ID in mapping")

	// Internal errors
	ErrInternal            = errors.New("bacnet internal error")
)
