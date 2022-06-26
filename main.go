package main

import (
	"os"
	"os/signal"
	"phanes/bll"
	"phanes/client"
	"phanes/collector"
	"phanes/config"
	"phanes/server"
	"phanes/store"
)

type InitFunc func() func()

func main() {
	var (
		exit    = make(chan struct{})
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
		close(exit)
	}()

	for _, fn := range bootstraps {
		cancels = append(cancels, fn())
	}

	<-exit
}
