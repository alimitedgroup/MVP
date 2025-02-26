package broker

import "github.com/nats-io/nats.go"

type NatsMessageBroker struct {
	nats *nats.Conn
}

func NewNatsMessageBroker() MessageBroker {
	nc, _ := nats.Connect(nats.DefaultURL)
	return &NatsMessageBroker{nc}
}
