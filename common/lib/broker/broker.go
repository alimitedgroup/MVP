package broker

import (
	"log"

	"github.com/nats-io/nats.go"
)

type MessageBroker interface {
}

type BrokerConfig struct {
	Url string `mapstructure:"url"`
}

type NatsMessageBroker struct {
	Nats *nats.Conn
	Js   nats.JetStream
}

func NewNatsMessageBroker(cfg *BrokerConfig) (*NatsMessageBroker, error) {
	log.Printf("Connecting to NATS at %s\n", cfg.Url)

	nc, err := nats.Connect(cfg.Url)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &NatsMessageBroker{nc, js}, nil
}

func (n *NatsMessageBroker) RequestSubscribe(subject Subject, queue Queue, handler RequestHandler) (*nats.Subscription, error) {
	sub, err := n.Nats.QueueSubscribe(subject.String(), queue.String(), func(msg *nats.Msg) {
		handler(msg)
	})
	if err != nil {
		return nil, err
	}

	return sub, nil
}
