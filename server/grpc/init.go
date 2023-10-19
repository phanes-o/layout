package grpc

import (
	"time"

	"github.com/asim/go-micro/plugins/registry/etcd/v4"
	"github.com/asim/go-micro/plugins/server/grpc/v4"
	"github.com/phanes-o/utils"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"phanes/config"
	"phanes/server/grpc/middleware"
)

func Init() micro.Option {

	opts := []server.Option{
		server.Name(config.Conf.Name + "-grpc"),
		server.Version(config.Conf.Version),
		server.RegisterTTL(time.Second * 30),
		server.RegisterInterval(time.Second * 15),
		server.Registry(etcd.NewRegistry(registry.Addrs(config.EtcdAddr))),
		server.WrapHandler(middleware.ServerTraceWrapper()),
		server.WrapHandler(middleware.Log()),
	}

	if config.Conf.Server.Grpc.GrpcListen != "" {
		opts = append(opts, server.Address(config.Conf.Server.Grpc.GrpcListen))
	}

	if config.Conf.Server.Grpc.DiscoveryListen != "" {
		opts = append(opts, server.Advertise(config.Conf.Server.Grpc.DiscoveryListen))
	}

	srv := grpc.NewServer(opts...)
	// register grpc services
	// example: utils.Throw(micro.RegisterHandler(srv, new(App)))

	// ⚠️Waring!!!: Your service struct Name Must seem to the .proto file service Name
	// utils.Throw(micro.RegisterHandler(srv, new(v1.User)))

	utils.Throw(srv.Start())
	return micro.Server(srv)
}
