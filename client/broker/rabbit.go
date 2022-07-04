package broker

import (
	"github.com/asim/go-micro/plugins/broker/rabbitmq/v4"
	"go-micro.dev/v4/broker"
	"phanes/config"
	"phanes/utils"
)

var defaultRabbitMQAddress = "amqp://guest:guest@localhost:5672"

func InitRabbit() broker.Broker {
	if config.Conf.Broker.Addr != "" {
		defaultRabbitMQAddress = config.Conf.Broker.Addr
	}
	b := rabbitmq.NewBroker(broker.Addrs(defaultRabbitMQAddress))
	utils.Throw(b.Connect())
	return b
}
