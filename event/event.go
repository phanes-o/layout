package event

import (
	"reflect"
	"sync"
)

// Handler 定义了事件处理函数的签名。
type Handler func(payload interface{})

// EventBus 是一个事件总线。
type EventBus struct {
	handlers map[Event][]Handler
	mutex    sync.RWMutex
}

// NewEventBus 创建一个新的事件总线。
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[Event][]Handler),
	}
}

// Register 在事件总线中注册一个事件处理函数。
func (bus *EventBus) Register(event Event, handler Handler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	handlers, ok := bus.handlers[event]
	if !ok {
		handlers = []Handler{}
	}

	bus.handlers[event] = append(handlers, handler)
}

// Unregister 在事件总线中取消注册一个事件处理函数。
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

// PublishAsync 发布一个异步事件。
func (bus *EventBus) PublishAsync(event Event, payload interface{}) {
	bus.mutex.RLock()
	handlers, ok := bus.handlers[event]
	bus.mutex.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		go handler(payload)
	}
}

// PublishSync 发布一个同步事件。
func (bus *EventBus) PublishSync(event Event, payload interface{}) {
	bus.mutex.RLock()
	handlers, ok := bus.handlers[event]
	bus.mutex.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		handler(payload)
	}
}
