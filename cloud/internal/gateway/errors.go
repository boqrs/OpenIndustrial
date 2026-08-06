package gateway

import "errors"

var (
	ErrGatewayNameRequired = errors.New("gateway name is required")
	ErrGatewayNotFound     = errors.New("gateway not found")
)