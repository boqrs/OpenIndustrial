package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type StatusMessageType string

const (
	StatusMessageHeartbeat StatusMessageType = "heartbeat"

	StatusMessageOnline StatusMessageType = "online"

	StatusMessageOffline StatusMessageType = "offline"

	StatusMessageCommandAck StatusMessageType = "command_ack"

	StatusMessageEvent StatusMessageType = "event"
)

// StatusMessage is the common Device -> Cloud message.
type StatusMessage struct {
	Type StatusMessageType `json:"type"`

	Timestamp *time.Time `json:"timestamp,omitempty"`

	CommandID uint `json:"command_id,omitempty"`

	Success *bool `json:"success,omitempty"`

	Error string `json:"error,omitempty"`

	Event string `json:"event,omitempty"`

	Data json.RawMessage `json:"data,omitempty"`
}

// ParseStatusMessage parses and validates a Device status message.
func ParseStatusMessage(
	payload []byte,
) (*StatusMessage, error) {
	if len(payload) == 0 {
		return nil, errors.New(
			"empty status message",
		)
	}

	var message StatusMessage

	if err := json.Unmarshal(
		payload,
		&message,
	); err != nil {
		return nil, fmt.Errorf(
			"invalid status message: %w",
			err,
		)
	}

	if err := message.Validate(); err != nil {
		return nil, err
	}

	return &message, nil
}

// Validate validates a Device status message.
func (m *StatusMessage) Validate() error {
	if m == nil {
		return errors.New(
			"status message is nil",
		)
	}

	switch m.Type {

	case StatusMessageHeartbeat,
		StatusMessageOnline,
		StatusMessageOffline:

		return nil

	case StatusMessageCommandAck:

		if m.CommandID == 0 {
			return errors.New(
				"command_id is required",
			)
		}

		if m.Success == nil {
			return errors.New(
				"success is required",
			)
		}

		return nil

	case StatusMessageEvent:

		if strings.TrimSpace(
			m.Event,
		) == "" {
			return errors.New(
				"event is required",
			)
		}

		return nil

	default:

		return fmt.Errorf(
			"unsupported status message type: %q",
			m.Type,
		)
	}
}
