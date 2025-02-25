package main

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/api_gateway/api"
	"github.com/alimitedgroup/MVP/srv/api_gateway/api/router"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

type NatsMessageBroker struct {
	nats *nats.Conn
}

type MessageBroker interface {
}

func NewNatsMessageBroker() MessageBroker {
	nc, _ := nats.Connect(nats.DefaultURL)
	return &NatsMessageBroker{nc}
}

func Run(h lib.HTTPHandler, routes router.Routes) {
	routes.Setup()

	_ = h.Engine.Run(":8080")
}

var Modules = fx.Options(
	fx.Provide(lib.Module),
	fx.Provide(api.Module),
	fx.Provide(NewNatsMessageBroker),
)

func main() {
	_ = context.Background()

	opts := fx.Options(Modules)
	app := fx.New(
		opts,
		fx.Invoke(Run),
	)

	app.Run()
}
