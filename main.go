package main

import (
	"context"
	"go-micro.dev/v4/logger"
	"os"
	"os/signal"
	"phanes/bll"
	"phanes/client"
	"phanes/collector"
	log "phanes/collector/logger"
	"phanes/config"
	"phanes/server"
	"phanes/store"
)

type InitFunc func() func()

func main() {
	var (
		cancels = make([]func(), 0)

		// system init
		bootstraps = []InitFunc{
			config.Init,
			collector.Init,
			store.Init,
			bll.Init,
			server.Init,
			client.Init,
		}
	)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, os.Kill)
		<-sigint

		for _, cancel := range cancels {
			cancel()
		}
		close(config.ExitC)
	}()

	for _, fn := range bootstraps {
		cancels = append(cancels, fn())
	}
	logger.Info("test go-micro logger")

	log.Info(context.Background(), "test log")

	<-config.ExitC
	log.Info(context.Background(), "server shutdown!")
}
