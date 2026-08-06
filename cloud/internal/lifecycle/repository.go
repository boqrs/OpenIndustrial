package lifecycle

import "context"

type Repository interface {
	SaveEvent(
		ctx context.Context,
		event *Event,
	) error

	GetEvents(
		ctx context.Context,
		productInstanceID string,
	) (
		[]Event,
		error,
	)
}