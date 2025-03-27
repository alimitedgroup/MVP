package controller

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type notificationRouterMessageBroker struct {
	nc *nats.Conn
	js jetstream.JetStream
}

func NewNotificationRouterMessageBroker(nc *nats.Conn, js jetstream.JetStream) notificationRouterMessageBroker {
	return notificationRouterMessageBroker{nc, js}
}
