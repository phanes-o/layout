package utils

import (
	"fmt"
	"net"
	"strings"

	"go-micro.dev/v4/server"
	"go-micro.dev/v4/util/addr"
	"phanes/config"
	"phanes/lib/traefik"
)

func Register(name string, srv server.Server) error {
	serverAddr, err := extractServerAddr(srv)
	Throw(err)
	var (
		reverseType traefik.ReverseType
		rulePrefix  string
	)

	for _, t := range config.Conf.Traefik.Type {
		switch t {
		case "http":
			reverseType = traefik.ReverseTypeHttp
			rulePrefix = "Host"
		case "grpc", "h2c":
			reverseType = traefik.ReverseTypeH2c
			rulePrefix = "Host"
		case "tcp":
			reverseType = traefik.ReverseTypeTcp
			rulePrefix = "HostSNI"
		case "udp":
			reverseType = traefik.ReverseTypeUdp
			rulePrefix = "HostSNI"
		}
		
		conf := &traefik.Config{
			Type:         reverseType,
			Tls:          config.Conf.Traefik.TLS,
			Enable:       config.Conf.Traefik.Enabled,
			SrvName:      name,
			SrvAddr:      serverAddr,
			Rule:         fmt.Sprintf("%s(`%s`) %s PathPrefix(`%s`)", rulePrefix, config.Conf.Traefik.Domain, config.Conf.Traefik.Rule, config.Conf.Traefik.Prefix),
			Prefix:       config.Conf.Traefik.Prefix,
			EndPoints:    []string{t},
			EnableRouter: config.Conf.Traefik.EnableRouter,
		}

		if config.Conf.Traefik.TLS {
			conf.EndPoints = []string{t}
		}
		if err := traefik.Register(conf); err != nil {
			return err
		}
	}
	return nil
}

func extractServerAddr(srv server.Server) (string, error) {
	host, port, err := net.SplitHostPort(srv.Options().Address)
	if err != nil {
		return "", err
	}
	extractAddr, err := addr.Extract(host)
	if err != nil {
		return "", err
	}
	if strings.Count(extractAddr, ":") > 0 {
		extractAddr = "[" + extractAddr + "]"
	}

	addr := fmt.Sprintf("%v:%v", extractAddr, port)
	return addr, nil
}
