package client

import (
	"phanes/client/broker"
	"phanes/client/grpc"
)

func Init() func() {
	var cancels = make([]func(), 0)
	var inits = []func() func(){
		grpc.Init,
		broker.Init,
	}

	for _, init := range inits {
		cancels = append(cancels, init())
	}
	
	return func() {
		for _, cancel := range cancels {
			cancel()
		}
	}
}
