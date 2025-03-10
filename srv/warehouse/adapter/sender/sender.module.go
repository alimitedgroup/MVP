package sender

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewPublishStockUpdateAdapter, fx.As(new(port.ICreateStockUpdatePort)))),
)
