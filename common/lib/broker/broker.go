package broker

import (
	"github.com/nats-io/nats.go"
)

type MessageBroker interface {
}

type BrokerConfig struct {
	Url string
}

type NatsMessageBroker struct {
	Nats *nats.Conn
	Js   nats.JetStream
}

func NewNatsMessageBroker(cfg *BrokerConfig) (*NatsMessageBroker, error) {
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
	sub, err := n.Nats.QueueSubscribe(subject.Name(), queue.String(), func(msg *nats.Msg) {
		handler(msg)
	})
	if err != nil {
		return nil, err
	}

	return sub, nil
}
