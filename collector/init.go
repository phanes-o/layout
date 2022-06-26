package collector

import "phanes/collector/logger"

func Init() func() {
	var (
		cancels = make([]func(), 0)

		inits = []func() func(){
			logger.Init,
		}
	)

	for _, init := range inits {
		cancels = append(cancels, init())
	}

	return func() {
		for _, cancel := range cancels {
			cancel()
		}
	}
}
