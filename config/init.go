package config

import (
	"github.com/asim/go-micro/plugins/config/source/etcd/v4"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4/cmd"
	"go-micro.dev/v4/config"
	"os"
	"phanes/utils"
	"time"
)

func Init() func() {
	var (
		conf config.Config
		err  error
	)
	utils.Throw(extractEtcdAddr())
	etcdSource := etcd.NewSource(
		etcd.WithAddress(EtcdAddr),
		etcd.WithPrefix(prefix),
		etcd.StripPrefix(true),
		etcd.WithDialTimeout(time.Second*time.Duration(10)),
	)
	// Create new config
	if conf, err = config.NewConfig(); err != nil {
		utils.Throw(err)
	}

	// Load file source
	if err = conf.Load(etcdSource); err != nil {
		utils.Throw(err)
	}
	err = conf.Scan(&Conf)
	// init others

	return func() {}
}

func extractEtcdAddr() error {
	app := cli.NewApp()
	app.Name = "vsm-manager"
	app.Flags = cmd.DefaultCmd.App().Flags
	app.Action = func(ctx *cli.Context) error {
		EtcdAddr = ctx.String("registry_address")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}
	return nil
}
