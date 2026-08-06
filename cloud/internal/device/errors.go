package device

import "errors"

var (
	ErrDeviceNameRequired = errors.New("device name is required")
	ErrDeviceNotFound     = errors.New("device not found")
)