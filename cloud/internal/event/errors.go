package event

import "errors"

var (
	ErrEventNameEmpty = errors.New("event name is empty")
	ErrHandlerNil     = errors.New("handler is nil")
)