package persistence

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewStockRepositoryIml, fx.As(new(StockRepository)))),
	fx.Provide(fx.Annotate(NewStockPersistanceAdapter,
		fx.As(new(port.SaveUpdateStockPort))),
	),
)
