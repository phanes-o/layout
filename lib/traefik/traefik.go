package traefik

import (
	"time"
)

type ReverseType int

const (
	ReverseTypeHttp ReverseType = iota
	ReverseTypeTcp
	ReverseTypeH2c
	ReverseTypeUdp
)

type Config struct {
	TTL time.Duration

	EnableRouter bool
	Tls          bool
	Enable       bool
	Type         ReverseType
	SrvName      string
	SrvAddr      string
	Rule         string
	Prefix       string
	EndPoints    []string
}

func Init(addr string) error {
	return etcd.init(addr)
}

func Register(conf *Config) error {
	return etcd.register(conf)
}
