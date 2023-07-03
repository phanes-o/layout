package bll

type Service interface {
	init() func()
}

// register service
var services = []Service{}

func Register(srv Service) {
	services = append(services, srv)
}

func Init() func() {
	var cancels = make([]func(), 0)

	for _, srv := range services {
		cancels = append(cancels, srv.init())
	}

	return func() {
		for _, cancel := range cancels {
			cancel()
		}
	}
}
