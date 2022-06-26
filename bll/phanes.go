package bll

import "phanes/event"

var Phanes = &phanes{}

type phanes struct {
}

func (p phanes) onEvent(ed *event.Data) {
}

func (p phanes) init() func() {

	return func() {
		
	}
}
