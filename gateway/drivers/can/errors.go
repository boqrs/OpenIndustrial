package can

import "errors"

var (
	ErrNotConnected      = errors.New("can: not connected")
	ErrAlreadyConnected  = errors.New("can: already connected")
	ErrInvalidInterface  = errors.New("can: invalid interface")
	ErrSignalNotFound    = errors.New("can: signal not found")
	ErrFrameNotFound     = errors.New("can: frame not found for signal")
	ErrEncodingFailed    = errors.New("can: failed to encode signal")
	ErrDecodingFailed    = errors.New("can: failed to decode frame")
)