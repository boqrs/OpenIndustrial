package workorder

import "errors"

var (
	ErrInvalidQuantity = errors.New("work order quantity must be positive")
	ErrWorkOrderNotFound = errors.New("work order not found")
)