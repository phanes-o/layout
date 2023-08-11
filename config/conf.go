package config

import (
	"go-micro.dev/v4"
)

var (
	Conf         *Config
	MicroService micro.Service
	EtcdAddr     = ""
	configFile   = ""
	ExitC        = make(chan bool)
	prefix       = "/phanes/config/hello"
)

type Config struct {
	Name        string `json:"name" yaml:"name" toml:"name"`
	Env         string `json:"env" yaml:"env" toml:"env"`
	Version     string `json:"version" yaml:"version" toml:"version"`
	AutoMigrate bool   `json:"auto_migrate" yaml:"auto_migrate" toml:"auto_migrate"`

	Http struct {
		Listen        string `json:"listen" yaml:"listen" toml:"listen"`
		ValidateTrans string `json:"validate_trans" yaml:"validateTrans" toml:"validate_trans"`
	}

	Collect Collect `json:"collect" yaml:"collect" toml:"collect"`

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
		Type    string `json:"type" yaml:"type" toml:"type"`          // support tcp,http,grpc,udp
		Rule    string `json:"rule" yaml:"rule" yaml:"toml"`          // router rule  value: "||", "&&"
		TLS     bool   `json:"tls" yaml:"tls" yaml:"tls"`             // is or not enable tls
		Enabled bool   `json:"enabled" yaml:"enabled" toml:"enabled"` // is or not register traefik
		Domain  string `json:"domain" yaml:"domain" toml:"domain"`    // gateway domain
		Prefix  string `json:"prefix" yaml:"prefix" toml:"prefix"`    // if Prefix is not empty, it will register router and middleware
	} `json:"traefik" yaml:"traefik" toml:"traefik"`
}

type Collect struct {
	Log struct {
		LogLevel   int8   `json:"log_level" yaml:"log_level" toml:"log_level"` // log level support -1:5
		Prefix     string `json:"prefix" yaml:"prefix" toml:"prefix"`
		FileName   string `json:"file_name"  yaml:"file_name" toml:"file_name"`
		BufferSize int    `json:"buffer_size" yaml:"buffer_size" toml:"buffer_size"`
		Interval   int64  `json:"interval" yaml:"interval" toml:"interval"`
	} `json:"log" yaml:"log" toml:"log"`

	Trace struct {
		Protocol string `json:"protocol" yaml:"protocol" toml:"protocol"` // trace report way "http" or "grpc"
		Type     string `json:"type" yaml:"type" toml:"type"`             // trace report type "otel" or "jaeger" or "zipkin"
		Addr     string `json:"addr" yaml:"addr" toml:"addr"`
	} `json:"trace" yaml:"trace" toml:"trace"`

	Metric struct {
		Listen string `json:"listen" yaml:"listen" toml:"listen"` // prometheus fatch addr
	} `json:"metric" yaml:"metric" toml:"metric"`
}
