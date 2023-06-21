package event

import (
	"reflect"
	"sync"
)

type Handler func(e Event, payload interface{})

type EventBus struct {
	handlers map[Event][]Handler
	mutex    sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[Event][]Handler),
	}
}

func (bus *EventBus) Register(event Event, handler Handler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	handlers, ok := bus.handlers[event]
	if !ok {
		handlers = []Handler{}
	}

	bus.handlers[event] = append(handlers, handler)
}

func (bus *EventBus) Unregister(event Event, handler Handler) {
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

func (bus *EventBus) PublishAsync(event Event, payload interface{}) {
	bus.mutex.RLock()
	handlers, ok := bus.handlers[event]
	bus.mutex.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		go handler(event, payload)
	}
}

func (bus *EventBus) PublishSync(event Event, payload interface{}) {
	bus.mutex.RLock()
	handlers, ok := bus.handlers[event]
	bus.mutex.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		handler(event, payload)
	}
}
