package controller

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type catalogRouterMessageBroker struct {
	nc *nats.Conn
	js jetstream.JetStream
}

func NewCatalogRouterMessageBroker(nc *nats.Conn, js jetstream.JetStream) catalogRouterMessageBroker {
	return catalogRouterMessageBroker{nc, js}
}
