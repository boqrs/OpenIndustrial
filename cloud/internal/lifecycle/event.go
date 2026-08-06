package lifecycle

import "time"

type Event struct {
	ID string

	ProductInstanceID string

	EventType string

	FromState string

	ToState string

	Source string

	Operator string

	Payload map[string]any

	CreatedAt time.Time
}