package broker

import (
	"go-micro.dev/v4/broker"
	"phanes/config"
)

var PubSub broker.Broker

func Init() func() {
	switch config.Conf.Broker.Type {
	case "rabbitmq":
		PubSub = InitRabbit()
	case "nats":
		PubSub = InitNats()
	}

	return func() {
		PubSub.Disconnect()
	}
}
