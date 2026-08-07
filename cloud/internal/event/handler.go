package event

// Handler handles domain events.
type Handler interface {
	Handle(
		event Event,
	) error
}