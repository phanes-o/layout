package broker

import (
	"github.com/asim/go-micro/plugins/broker/nats/v4"
	"github.com/phanes-o/utils"
	"go-micro.dev/v4/broker"
	"phanes/config"
)

var defaultNatsAddress = "nats://127.0.0.1:1222"

func InitNats() broker.Broker {
	if config.Conf.Client.Broker.Addr != "" {
		defaultNatsAddress = config.Conf.Client.Broker.Addr
	}
	b := nats.NewBroker(broker.Addrs(defaultNatsAddress))
	utils.Throw(b.Connect())
	return b
}
