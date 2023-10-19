package broker

import (
	"go-micro.dev/v4/broker"
	log "phanes/collector/logger"
	"phanes/config"
	"phanes/errors"
)

var pubSub broker.Broker

func Init() func() {
	switch config.Conf.Client.Broker.Type {
	case "rabbitmq":
		pubSub = InitRabbit()
	case "nats":
		pubSub = InitNats()
	}

	return func() {
		if pubSub != nil {
			if err := pubSub.Disconnect(); err != nil {
				log.Error(err.Error())
			}
		}
	}
}

func Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {
	if pubSub == nil {
		return errors.New("broker not init")
	}
	return pubSub.Publish(topic, m, opts...)
}

func Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	if pubSub == nil {
		return nil, errors.New("broker not init")
	}
	return pubSub.Subscribe(topic, handler, opts...)
}
