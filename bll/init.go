package bll

import "phanes/event"

type Service interface {
	onEvent(ed *event.Data)
	init() func()
}

var services = []Service{
	User,
}

func Init() func() {
	var cancels = make([]func(), 0)

	for _, srv := range services {
		cancels = append(cancels, srv.init())
	}

	event.Register(func(data *event.Data) {
		for _, srv := range services {
			srv.onEvent(data)
		}
	})

	return func() {
		for _, cancel := range cancels {
			cancel()
		}
	}
}
