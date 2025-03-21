package sender

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewPublishStockUpdateAdapter, fx.As(new(port.ICreateStockUpdatePort)))),
	fx.Provide(fx.Annotate(NewPublishReservationEventAdapter, fx.As(new(port.IStoreReservationEventPort)))),
)
