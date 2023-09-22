package event

import (
	"context"

	log "phanes/collector/logger"
)

var bus Bus

func Init() func() {
	bus = newEventBus()
	return func() {}
}

type Event string

type Handler func(ctx context.Context, e Event, payload interface{})

type Bus interface {
	subscribe(event Event, handler Handler)
	unsubscribe(event Event, handler Handler)
	publish(ctx context.Context, event Event, payload interface{})
}

func Subscribe(event Event, handler Handler) {
	if bus == nil {
		log.Error("event bus not init")
		return
	}
	bus.subscribe(event, handler)
}

func Unsubscribe(event Event, handler Handler) {
	if bus == nil {
		log.Error("event bus not init")
		return
	}
	bus.unsubscribe(event, handler)
}

func Publish(ctx context.Context, event Event, payload interface{}) {
	if bus == nil {
		log.Error("event bus not init")
		return
	}
	bus.publish(ctx, event, payload)
}
