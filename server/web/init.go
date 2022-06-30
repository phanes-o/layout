package web

import (
	"fmt"
	"github.com/asim/go-micro/plugins/registry/etcd/v4"
	"github.com/asim/go-micro/plugins/server/http/v4"
	"phanes/model"
	"phanes/server/web/middleware"
	"phanes/server/web/v1"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go-micro.dev/v4/util/addr"
	"net"
	"phanes/config"
	"phanes/lib/traefik"
	"phanes/utils"
	"strings"
	"time"
)

var (
	webName           string
	webAddr           string
	defaultListenAddr = ":7701"
	srv               server.Server
)

func Init() micro.Option {
	webName = config.Conf.Name + "-http"

	if config.Conf.HttpListen != "" {
		defaultListenAddr = config.Conf.HttpListen
	}

	srv = http.NewServer(
		server.Name(webName),
		server.Version(config.Conf.Version),
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
		utils.Throw(Register())
	}
	return micro.Server(srv)
}

func extractWebAddr(srv server.Server) error {
	host, port, err := net.SplitHostPort(srv.Options().Address)
	if err != nil {
		return err
	}
	extractAddr, err := addr.Extract(host)
	if err != nil {
		return err
	}
	if strings.Count(extractAddr, ":") > 0 {
		extractAddr = "[" + extractAddr + "]"
	}

	webAddr = fmt.Sprintf("%v:%v", extractAddr, port)
	return nil
}

func Register() error {
	utils.Throw(extractWebAddr(srv))

	conf := &traefik.Config{
		SrvName:   webName,
		SrvAddr:   webAddr,
		Rule:      fmt.Sprintf("Host(`%s`) || PathPrefix(`%s`)", config.Conf.Traefik.Domain, config.Conf.Traefik.Prefix),
		Prefix:    config.Conf.Traefik.Prefix,
		EndPoints: []string{"http"},
	}
	if strings.ToLower(config.Conf.Env) == model.EnvProd {
		conf.EndPoints = []string{"https"}
	}
	return traefik.Register(conf)
}
