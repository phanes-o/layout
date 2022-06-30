package config

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	KV *redis.Client
	DB *gorm.DB

	Conf     = &Config{}
	EtcdAddr = ""
	ExitC    = make(chan bool)
	prefix   = "/phanes/config/hello"
)

type Config struct {
	Name    string `json:"name"`
	Env     string `json:"env"`
	Version string `json:"version"`

	HttpListen string `json:"http_listen"`
	Jaeger     string `json:"jaeger"`

	Postgres string `json:"postgres"`

	Log struct {
		FileName string `json:"file_name"`
		RedisKey string `json:"redis_key"`
	} `json:"log"`

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
