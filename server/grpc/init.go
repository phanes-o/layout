package grpc

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v4"
	"github.com/asim/go-micro/plugins/server/grpc/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"phanes/config"
	"phanes/server/grpc/middleware"
	v1 "phanes/server/grpc/v1"
	"phanes/utils"
	"time"
)

func Init() micro.Option {
	srv := grpc.NewServer(
		server.Name(config.Conf.Name+"-grpc"),
		server.Version(config.Conf.Version),
		server.RegisterTTL(time.Second*30),
		server.RegisterInterval(time.Second*15),
		server.Registry(etcd.NewRegistry(registry.Addrs(config.EtcdAddr))),
		server.WrapHandler(middleware.Log()),
	)
	// register grpc services
	// example: utils.Throw(micro.RegisterHandler(srv, new(App)))
	utils.Throw(micro.RegisterHandler(srv, new(v1.User)))

	utils.Throw(srv.Start())
	return micro.Server(srv)
}
