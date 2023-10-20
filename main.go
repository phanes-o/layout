package main

import (
	"context"

	"phanes/assistant"
	"phanes/bll"
	"phanes/cache"
	"phanes/client"
	"phanes/collector"
	"phanes/config"
	"phanes/event"
	"phanes/server"
	"phanes/store"
)

type InitFunc func() func()

func main() {
	var (
		// system init func
		bootstraps = []InitFunc{
			config.Init,
			collector.Init,
			cache.Init,
			server.Init,
			client.Init,
			store.Init,
			assistant.Init,
			bll.Init,
		}
	)

	event.Init()
	for _, bootstrap := range bootstraps {
		cancel := bootstrap()
		event.Subscribe(event.EventExit, func(ctx context.Context, e event.Event, payload interface{}) {
			cancel()
		})
	}
	config.MicroService.Run()
}
