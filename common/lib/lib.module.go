package lib

import (
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewHTTPHandler),
	fx.Provide(broker.NewNatsConn),
	fx.Provide(broker.NewNatsMessageBroker),
	fx.Provide(broker.NewRestoreStreamControl),
	fx.Provide(broker.NewRestoreStreamControlFactory),
)
