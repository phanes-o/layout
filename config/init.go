package config

import (
	"log"
	"os"
	"time"

	"github.com/asim/go-micro/plugins/config/encoder/toml/v4"
	"github.com/asim/go-micro/plugins/config/encoder/yaml/v4"
	"github.com/asim/go-micro/plugins/config/source/etcd/v4"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4/cmd"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/encoder"
	ejson "go-micro.dev/v4/config/encoder/json"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source"
	"go-micro.dev/v4/config/source/file"
	"phanes/lib/traefik"
	"phanes/model"
	"phanes/utils"
)

func Init() func() {
	var (
		e              encoder.Encoder
		err            error
		conf           config.Config
		cancels        = make([]func(), 0)
		iSource        source.Source
		configFileType = model.ConfigFileTypeYaml
	)

	switch configFileType {
	case model.ConfigFileTypeJson:
		e = ejson.NewEncoder()
	case model.ConfigFileTypeYaml:
		e = yaml.NewEncoder()
	case model.ConfigFileTypeToml:
		e = toml.NewEncoder()
	}

	Conf = &Config{}
	utils.Throw(extractEtcdAddr())

	if configFile != "" {
		iSource = file.NewSource(file.WithPath(configFile))
	} else {
		iSource = etcd.NewSource(
			etcd.WithAddress(EtcdAddr),
			etcd.WithPrefix(prefix),
			source.WithEncoder(e),
			etcd.StripPrefix(true),
			etcd.WithDialTimeout(time.Second*time.Duration(10)),
		)
	}
	// Create new config
	conf, err = config.NewConfig(
		config.WithSource(iSource),
		config.WithReader(json.NewReader(reader.WithEncoder(e))),
	)
	utils.Throw(err)
	utils.Throw(conf.Scan(Conf))

	watcher, err := conf.Watch()
	if err != nil {
		log.Println(err)
	}

	go func() {
		for {
			select {
			case <-ExitC:
				watcher.Stop()
				return
			default:
			}
			next, err := watcher.Next()
			if err = next.Scan(&Conf); err != nil {
				log.Println(err)
			}
		}
	}()

	utils.Throw(traefik.Init(EtcdAddr))
	// init others
	inits := []func() func(){}
	for _, init := range inits {
		cancels = append(cancels, init())
	}
	return func() {
		for _, cancel := range cancels {
			cancel()
		}
	}
}

func extractEtcdAddr() error {
	app := cli.NewApp()
	app.Name = "phanes-layout"
	app.Flags = cmd.DefaultCmd.App().Flags
	app.Action = func(ctx *cli.Context) error {
		EtcdAddr = ctx.String("registry_address")
		configFile = ctx.String("config")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}
	return nil
}
