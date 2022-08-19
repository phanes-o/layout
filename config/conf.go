package config

import (
	"github.com/go-redis/redis/v8"
	"go-micro.dev/v4"
)

var (
	KV           *redis.Client
	Conf         = &Config{}
	MicroService micro.Service
	EtcdAddr     = ""
	ExitC        = make(chan bool)
	prefix       = "/phanes/config/hello"
)

type Config struct {
	Name    string `json:"name"`
	Env     string `json:"env"`
	Version string `json:"version"`

	HttpListen string `json:"http_listen"`
	Jaeger     string `json:"jaeger"`

	Collect struct {
		Log struct {
			FileName string `json:"file_name"`
			RedisKey string `json:"redis_key"`
		} `json:"log"`
		Trace struct {
			Addr string `json:"addr"`
		} `json:"trace"`
		Metric struct {
			Addr string `json:"addr"`
		} `json:"metric"`
	} `json:"collect"`

	DB []struct {
		Addr string `json:"addr"` // host=127.0.0.1 user=root password=root dbname=signal port=5432 sslmode=disable TimeZone=Asia/Shanghai
		Type string `json:"type"` // postgres, mysql, sqlite, mongo
		User string `json:"user"` // if addr not like Addr example or other need, you should set
		Pwd  string `json:"pwd"`
	} `json:"db"`

	Redis struct {
		Addr string `json:"addr"`
		Pwd  string `json:"pwd"`
	} `json:"redis"`

	Broker struct {
		Type string `json:"type"` // support: rabbitmq, nats
		Addr string `json:"addr"`
		User string `json:"user"`
		Pwd  string `json:"pwd"`
	} `json:"broker"`

	Traefik struct {
		Enabled bool   `json:"enabled"`
		Domain  string `json:"domain"`
		Prefix  string `json:"prefix"`
	} `json:"traefik"`
}
