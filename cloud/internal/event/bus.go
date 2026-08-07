package event

type Bus interface {
	Publish(event Event) error
	Subscribe(eventName string, handler Handler) error
}