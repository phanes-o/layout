package bll

import (
	"phanes/event"
)

type eventHandle struct {
}

func (e *eventHandle) init() func() {
	// register event handler
	event.Bus.Register(event.ExampleEvent, e.ExampleEventHandler)

	return func() {}
}

func (e *eventHandle) ExampleEventHandler(data interface{}) {
	// handle event
}
