package persistence

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewStockRepositoryIml, fx.As(new(StockRepository)))),
	fx.Provide(fx.Annotate(NewStockPersistanceAdapter, fx.As(new(port.GetStockPort)), fx.As(new(port.ApplyStockUpdatePort)))),
	fx.Provide(fx.Annotate(NewCatalogRepositoryIml, fx.As(new(CatalogRepository)))),
	fx.Provide(fx.Annotate(NewCatalogPersistanceAdapter, fx.As(new(port.ApplyCatalogUpdatePort)), fx.As(new(port.GetGoodPort)))),
)
