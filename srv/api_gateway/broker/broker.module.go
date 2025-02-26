package broker

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewNatsMessageBroker),
)

type MessageBroker interface {
}
