package trace

import (
	"context"
	"fmt"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/lifecycle"
)

type Service struct {
	lifecycleRepo lifecycle.Repository
}

func NewService(
	lifecycleRepo lifecycle.Repository,
) *Service {
	return &Service{
		lifecycleRepo: lifecycleRepo,
	}
}

func (s *Service) Query(
	ctx context.Context,
	productInstanceID string,
) ([]TimelineItem, error) {
	events, err := s.lifecycleRepo.GetEvents(
		ctx,
		productInstanceID,
	)
	if err != nil {
		return nil, err
	}

	var timeline []TimelineItem
	for _, event := range events {
		timeline = append(
			timeline,
			TimelineItem{
				Time: event.CreatedAt.Format("2006-01-02 15:04:05"),
				Type: event.EventType,
				Description: fmt.Sprintf(
					"%s -> %s",
					event.FromState,
					event.ToState,
				),
			},
		)
	}

	return timeline, nil
}