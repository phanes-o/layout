package event

type HandleFunc func(data *Data)

var listeners = make([]HandleFunc, 0, 0)

func Register(fn HandleFunc) {
	if fn == nil {
		return
	}
	listeners = append(listeners, fn)
}

func Emit(data *Data) {
	if data == nil {
		return
	}

	for _, fn := range listeners {
		if fn != nil {
			fn(data)
		}
	}
}
