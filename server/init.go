package server

import (
	"github.com/asim/go-micro/plugins/client/grpc/v4"
	"github.com/asim/go-micro/plugins/registry/etcd/v4"
	"github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"phanes/config"
	grpcServer "phanes/server/grpc"

	webServer "phanes/server/web"
	// example: other server
	// exampleServer "phanes/server/example_server"
)

var Service micro.Service

func Init() func() {
	Service = micro.NewService()

	Service.Init(
		micro.Registry(etcd.NewRegistry(registry.Addrs(config.EtcdAddr))),
		micro.AfterStop(AfterExit),
		micro.AfterStop(AfterExit),
		micro.Client(grpc.NewClient()),
		// choose you needed wrapper
		micro.WrapClient(roundrobin.NewClientWrapper()),
		grpcServer.Init(),
		webServer.Init(),
		// Add other server here
		// example:
		// exampleServer.Init(),
	)

	return func() {}
}

func AfterExit() error {

	return nil
}
