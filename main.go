package main

import (
	"os"
	"os/signal"

	sig "go-micro.dev/v4/util/signal"
	"phanes/assistant"
	"phanes/bll"
	"phanes/client"
	"phanes/collector"
	log "phanes/collector/logger"
	"phanes/config"
	"phanes/event"
	"phanes/server"
	"phanes/store"
)

type InitFunc func() func()

func main() {
	var (
		cancels = make([]func(), 0)

		// system init func
		bootstraps = []InitFunc{
			config.Init,
			collector.Init,
			event.Init,
			server.Init,
			client.Init,
			store.Init,
			assistant.Init,
			bll.Init,
		}
	)
	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, sig.Shutdown()...)
		<-sigint

		for _, cancel := range cancels {
			cancel()
		}
		close(config.ExitC)
	}()

	for _, fn := range bootstraps {
		cancels = append(cancels, fn())
	}
	log.Info("finished to init all component")

	<-config.ExitC
	log.Info("server shutdown!")
}
