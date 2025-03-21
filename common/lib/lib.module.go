package lib

import (
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(broker.NewNatsMessageBroker),
	fx.Provide(broker.NewRestoreStreamControl),
	fx.Provide(fx.Annotate(broker.NewRestoreStreamControlFactory, fx.As(new(broker.IRestoreStreamControlFactory)))),
)
