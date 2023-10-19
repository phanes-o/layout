package broker

import (
	"github.com/asim/go-micro/plugins/broker/rabbitmq/v4"
	"github.com/phanes-o/utils"
	"go-micro.dev/v4/broker"
	"phanes/config"
)

var defaultRabbitMQAddress = "amqp://guest:guest@localhost:5672"

func InitRabbit() broker.Broker {
	if !config.Conf.Client.Broker.Enabled {
		return nil
	}
	if config.Conf.Client.Broker.Addr != "" {
		defaultRabbitMQAddress = config.Conf.Client.Broker.Addr
	}
	b := rabbitmq.NewBroker(broker.Addrs(defaultRabbitMQAddress))
	utils.Throw(b.Connect())
	return b
}
