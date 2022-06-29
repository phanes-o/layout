package broker

import (
	"github.com/asim/go-micro/plugins/broker/rabbitmq/v4"
	"go-micro.dev/v4/broker"
	"phanes/config"
)

var defaultRabbitMQAddress = "amqp://guest:guest@localhost:5672"

func InitRabbit() broker.Broker {
	if config.Conf.Broker.Addr != "" {
		defaultRabbitMQAddress = config.Conf.Broker.Addr
	}
	return rabbitmq.NewBroker(broker.Addrs(defaultRabbitMQAddress))
}
