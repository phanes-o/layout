package config

import (
	"github.com/go-redis/redis/v8"
	"go-micro.dev/v4"
)

var (
	KV           *redis.Client
	Conf         *Config
	MicroService micro.Service
	EtcdAddr     = ""
	configFile   = ""
	ExitC        = make(chan bool)
	prefix       = "/phanes/config/hello"
)

type Config struct {
	Base struct {
		Name       string `json:"name" yaml:"name" toml:"name"`
		Env        string `json:"env" yaml:"env" toml:"env"`
		Version    string `json:"version" yaml:"version" toml:"version"`
		HttpListen string `json:"http_listen" yaml:"http_listen" toml:"http_listen"`
	} `json:"base" yaml:"base" toml:"base"`

	Collect struct {
		Log struct {
			LogLevel uint8  `json:"log_level" yaml:"log_level" toml:"log_level"` // log level support -1:5
			Prefix   string `json:"prefix" yaml:"prefix" toml:"prefix"`
			FileName string `json:"file_name"  yaml:"file_name" toml:"file_name"`
			Redis    struct {
				RedisKey string `json:"redis_key" yaml:"redis_key" toml:"redis_key`
				Addr     string `json:"addr" yaml:"addr" toml:"addr"`
				Pwd      string `json:"pwd" yaml:"pwd" toml:"pwd"`
			} `json:"redis" yaml:"redis" toml:"redis"`
		} `json:"log" yaml:"log" toml:"log"`
		Trace struct {
			Addr string `json:"addr" yaml:"addr" toml:"addr"`
		} `json:"trace" yaml:"trace" toml:"trace"`
		Metric struct {
			Addr string `json:"addr" yaml:"addr" toml:"addr"`
		} `json:"metric" yaml:"metric" toml:"metric"`
	} `json:"collect" yaml:"collect" toml:"collect"`

	DB []struct {
		Addr string `json:"addr" yaml:"addr" toml:"addr"` // host=127.0.0.1 user=root password=root dbname=signal port=5432 sslmode=disable TimeZone=Asia/Shanghai
		Type string `json:"type" yaml:"type" toml:"type"` // postgres, mysql, sqlite, mongo
		User string `json:"user" yaml:"user" toml:"user"` // if addr not like Addr example or other need, you should set
		Pwd  string `json:"pwd" yaml:"pwd" toml:"pwd"`
	} `json:"db" yaml:"db" toml:"db"`

	Broker struct {
		Type string `json:"type" yaml:"type" toml:"type"` // support: rabbitmq, nats
		Addr string `json:"addr" yaml:"addr" toml:"addr"`
		User string `json:"user" yaml:"user" toml:"user"`
		Pwd  string `json:"pwd" yaml:"pwd" toml:"pwd"`
	} `json:"broker" yaml:"broker" toml:"broker"`

	Traefik struct {
		EnableRouter bool     `json:"enable_router" yaml:"enable_router" toml:"enable_router"`
		Type         []string `json:"type" yaml:"type", toml:"type"`
		// router rule  value: "||", "&&"
		Rule    string `json:"rule" yaml:"rule" yaml:"toml"`
		TLS     bool   `json:"tls" yaml:"tls" yaml:"tls"`
		Enabled bool   `json:"enabled" yaml:"enabled" toml:"enabled"`
		Domain  string `json:"domain" yaml:"domain" toml:"domain"`
		Prefix  string `json:"prefix" yaml:"prefix" toml:"prefix"`
	} `json:"traefik" yaml:"traefik" toml:"traefik"`
}
