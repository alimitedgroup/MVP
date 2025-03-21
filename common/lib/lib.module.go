package lib

import (
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(observability.New),
	fx.Provide(broker.NewNatsMessageBroker),
	fx.Provide(broker.NewRestoreStreamControl),
	fx.Provide(fx.Annotate(broker.NewRestoreStreamControlFactory, fx.As(new(broker.IRestoreStreamControlFactory)))),
)
