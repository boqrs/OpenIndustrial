package repository

import (
	"context"
	"encoding/json"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/lifecycle"
	"gorm.io/gorm"
)

type LifecycleRepository struct {
	db *gorm.DB
}

func NewLifecycleRepository(db *gorm.DB) *LifecycleRepository {
	return &LifecycleRepository{
		db: db,
	}
}

func (r *LifecycleRepository) SaveEvent(
	ctx context.Context,
	event *lifecycle.Event,
) error {
	payload, _ := json.Marshal(event.Payload)

	model := LifecycleEventModel{
		ID:                event.ID,
		ProductInstanceID: event.ProductInstanceID,
		EventType:         event.EventType,
		FromState:         event.FromState,
		ToState:           event.ToState,
		Source:            event.Source,
		Payload:           payload,
		CreatedAt:         event.CreatedAt,
	}

	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *LifecycleRepository) GetEvents(
	ctx context.Context,
	productInstanceID string,
) ([]lifecycle.Event, error) {
	var models []LifecycleEventModel
	err := r.db.WithContext(ctx).
		Where("product_instance_id = ?", productInstanceID).
		Order("created_at asc").
		Find(&models).
		Error
	if err != nil {
		return nil, err
	}

	var events []lifecycle.Event
	for _, model := range models {
		var payload map[string]any
		if model.Payload != nil {
			_ = json.Unmarshal(model.Payload, &payload)
		}

		events = append(events, lifecycle.Event{
			ID:                model.ID,
			ProductInstanceID: model.ProductInstanceID,
			EventType:         model.EventType,
			FromState:         model.FromState,
			ToState:           model.ToState,
			Source:            model.Source,
			Payload:           payload,
			CreatedAt:         model.CreatedAt,
		})
	}

	return events, nil
}