package event

import (
	"context"
	"reflect"
	"sync"
)

type eventBus struct {
	handlers map[Event][]Handler
	mutex    sync.RWMutex
}

func newEventBus() Bus {
	return &eventBus{
		handlers: make(map[Event][]Handler),
	}
}

func (bus *eventBus) subscribe(event Event, handler Handler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	handlers, ok := bus.handlers[event]
	if !ok {
		handlers = []Handler{}
	}

	bus.handlers[event] = append(handlers, handler)
}

func (bus *eventBus) unsubscribe(event Event, handler Handler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	handlers, ok := bus.handlers[event]
	if !ok {
		return
	}

	for i, h := range handlers {
		if reflect.DeepEqual(h, handler) {
			bus.handlers[event] = append(handlers[:i], handlers[i+1:]...)
			return
		}
	}
}

func (bus *eventBus) publish(ctx context.Context, event Event, payload interface{}) {
	bus.mutex.RLock()
	handlers, ok := bus.handlers[event]
	bus.mutex.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		go handler(ctx, event, payload)
	}
}
