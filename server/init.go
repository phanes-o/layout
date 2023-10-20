package server

import (
	"context"

	"github.com/asim/go-micro/plugins/client/grpc/v4"
	"github.com/asim/go-micro/plugins/registry/etcd/v4"
	"github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	log "phanes/collector/logger"
	"phanes/config"
	"phanes/event"
	grpcServer "phanes/server/grpc"
	"phanes/server/grpc/middleware"

	webServer "phanes/server/web"
	// example: other server
	// exampleServer "phanes-layout/server/example_server"
)

func Init() func() {
	config.MicroService = micro.NewService()

	config.MicroService.Init(
		micro.Registry(etcd.NewRegistry(registry.Addrs(config.EtcdAddr))),
		micro.AfterStart(AfterStart),
		micro.AfterStop(AfterExit),
		// client trace wrapper
		micro.Client(grpc.NewClient(client.WrapCall(middleware.ClientTraceWrapper()), client.Retries(0))),
		// choose you needed wrapper
		micro.WrapClient(roundrobin.NewClientWrapper()),
		webServer.Init(),
		grpcServer.Init(),
		// Add other server here
		// example:
		// exampleServer.Init(),
	)

	return func() {}
}

func AfterStart() error {
	log.Info("finished to init all component")
	return nil
}

func AfterExit() error {
	event.Publish(context.Background(), event.EventExit, nil)
	log.Info("server shutdown!")
	return nil
}
