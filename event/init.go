package event

var Bus *EventBus

func Init() func() {
	Bus = NewEventBus()
	return func() {}
}
