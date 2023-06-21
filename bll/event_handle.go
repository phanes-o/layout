package bll

import (
	"phanes/event"
)

type eventHandle struct {
}

func (h *eventHandle) init() func() {
	// register event handler
	event.Bus.Register(event.ExampleEvent, h.ExampleEventHandler)

	return func() {}
}

func (h *eventHandle) ExampleEventHandler(e event.Event, data interface{}) {
	// handle event
}
