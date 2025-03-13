package persistence

import (
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewStockRepositoryImpl, fx.As(new(IStockRepository)))),
	fx.Provide(fx.Annotate(NewOrderPersistanceAdapter, fx.As(new(port.IGetOrderPort)))),
	fx.Provide(fx.Annotate(NewStockPersistanceAdapter, fx.As(new(port.IGetStockPort)), fx.As(new(port.IApplyStockUpdatePort)))),
)
