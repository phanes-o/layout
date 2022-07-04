package collector

import (
	"phanes/collector/logger"
	"phanes/collector/trace"
)

func Init() func() {
	var (
		cancels = make([]func(), 0)

		inits = []func() func(){
			trace.Init,
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
