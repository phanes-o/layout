package config

import (
	"github.com/asim/go-micro/plugins/config/source/etcd/v4"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4/cmd"
	"go-micro.dev/v4/config"
	"log"
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
	if conf, err = config.NewConfig(config.WithSource(etcdSource)); err != nil {
		utils.Throw(err)
	}
	err = conf.Scan(&Conf)

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

	// init others
	initRedis()
	return func() {}
}

func extractEtcdAddr() error {
	app := cli.NewApp()
	app.Name = "phanes"
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
