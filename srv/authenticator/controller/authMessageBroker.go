package controller

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type authRouterMessageBroker struct {
	nc *nats.Conn
	js jetstream.JetStream
}

func NewAuthRouterMessageBroker(nc *nats.Conn, js jetstream.JetStream) authRouterMessageBroker {
	return authRouterMessageBroker{nc, js}
}
