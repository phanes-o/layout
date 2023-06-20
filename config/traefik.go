package config

import (
	"fmt"
	"net"
	"strings"

	"go-micro.dev/v4/server"
	"go-micro.dev/v4/util/addr"
	"phanes/lib/traefik"
	"phanes/utils"
)

func Register(name string, srv server.Server) error {
	serverAddr, err := extractServerAddr(srv)
	utils.Throw(err)
	var (
		reverseType traefik.ReverseType
		rulePrefix  string
	)

	switch Conf.Traefik.Type {
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
		Type:      reverseType,
		Tls:       Conf.Traefik.TLS,
		Enable:    Conf.Traefik.Enabled,
		SrvName:   name,
		SrvAddr:   serverAddr,
		Rule:      fmt.Sprintf("%s(`%s`) %s PathPrefix(`%s`)", rulePrefix, Conf.Traefik.Domain, Conf.Traefik.Rule, Conf.Traefik.Prefix),
		Prefix:    Conf.Traefik.Prefix,
		EndPoints: []string{Conf.Traefik.Type},
	}

	if Conf.Traefik.Type == "tcp" || Conf.Traefik.Type == "udp" {
		conf.Rule = fmt.Sprintf("%s(`%s`)", rulePrefix, Conf.Traefik.Domain)
	}
	if err := traefik.Register(conf); err != nil {
		return err
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
