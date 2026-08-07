package product

import "errors"

var (
	ErrModelNotFound      = errors.New("product model not found")
	ErrInstanceNotFound   = errors.New("product instance not found")
	ErrDuplicateSN        = errors.New("product instance with this SN already exists")
	ErrWorkOrderNotFound  = errors.New("work order not found")
	ErrInvalidEvent       = errors.New("invalid production event")
)