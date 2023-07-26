package assistant

var assistants = []Assistant{}

func Init() func() {
	var cancels = make([]func(), 0)

	for _, a := range assistants {
		cancels = append(cancels, a.Init())
	}

	return func() {
		for _, fn := range cancels {
			fn()
		}
	}
}

func Register(fn Assistant) {
	assistants = append(assistants, fn)
}

type Assistant interface {
	Init() func()
}
