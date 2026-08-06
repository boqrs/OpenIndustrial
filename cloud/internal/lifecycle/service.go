package lifecycle

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(
	repo Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}

func validateTransition(
	from string,
	to string,
	event string,
) bool {
	for _, t := range DefaultTransitions {
		if t.From == from &&
			t.To == to &&
			t.EventType == event {
			return true
		}
	}
	return false
}

func (s *Service) AppendEvent(
	ctx context.Context,
	productInstanceID string,
	from string,
	to string,
	eventType string,
	source string,
	payload map[string]any,
) error {
	if !validateTransition(
		from,
		to,
		eventType,
	) {
		return errors.New(
			"invalid lifecycle transition",
		)
	}
	event := &Event{
		ID:                uuid.New().String(),
		ProductInstanceID: productInstanceID,
		FromState:         from,
		ToState:           to,
		EventType:         eventType,
		Source:            source,
		Payload:           payload,
		CreatedAt:         time.Now(),
	}

	return s.repo.SaveEvent(
		ctx,
		event,
	)
}