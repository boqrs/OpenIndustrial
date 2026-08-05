package opcua

import "errors"

var (
	ErrNotConnected    = errors.New("opcua: not connected")
	ErrAlreadyConnected = errors.New("opcua: already connected")
	ErrSubscriptionFailed = errors.New("opcua: subscription failed")
)