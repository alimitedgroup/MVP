package broker

import (
	"github.com/nats-io/nats.go"
)

type MessageBroker interface {
}

type NatsMessageBroker struct {
	Nats *nats.Conn
	Js   nats.JetStream
}

func NewNatsMessageBroker() (*NatsMessageBroker, error) {
	nc, err := nats.Connect(nats.DefaultURL)
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
