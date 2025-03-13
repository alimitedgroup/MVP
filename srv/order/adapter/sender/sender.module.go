package sender

import (
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewPublishOrderUpdateAdapter, fx.As(new(port.ISaveOrderUpdatePort)))),
)
