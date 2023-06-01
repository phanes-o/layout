package web

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v4"
	"github.com/asim/go-micro/plugins/server/http/v4"
	"phanes/server/web/middleware"
	"phanes/server/web/v1"

	"time"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"phanes/config"
	"phanes/utils"
)

var (
	webName           string
	webAddr           string
	defaultListenAddr = ":7771"
	srv               server.Server
)

func Init() micro.Option {
	webName = config.Conf.Base.Name + "-http"

	if config.Conf.Base.HttpListen != "" {
		defaultListenAddr = config.Conf.Base.HttpListen
	}

	srv = http.NewServer(
		server.Name(webName),
		server.Version(config.Conf.Base.Version),
		server.Address(defaultListenAddr),
		server.RegisterTTL(time.Second*30),
		server.RegisterInterval(time.Second*15),
		server.Registry(etcd.NewRegistry(registry.Addrs(config.EtcdAddr))),
	)

	router := gin.New()
	gin.SetMode(gin.DebugMode)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// register routers
	v1Group := router.Group("v1", middleware.Log())
	v1.Init(v1Group)

	utils.Throw(srv.Handle(srv.NewHandler(router)))
	utils.Throw(srv.Start())
	if config.Conf.Traefik.Enabled {
		utils.Throw(config.Register(webName, srv))
	}
	return micro.Server(srv)
}
