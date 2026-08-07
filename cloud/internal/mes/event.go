package mes

import (
	"encoding/json"
	"time"
)

// Event represents a significant occurrence in the MES runtime.
type Event struct {
	ID                string
	Type              EventType
	ProductInstanceID string
	StationID         string
	TaskID            string
	Data              json.RawMessage
	Timestamp         time.Time
}

type EventType string

const (
	EventTaskCreated  EventType = "task.created"
	EventTaskStarted  EventType = "task.started"
	EventTaskFinished EventType = "task.finished" // Includes pass/fail data
	EventProductEnter EventType = "product.station.enter"
	EventProductExit  EventType = "product.station.exit"
)