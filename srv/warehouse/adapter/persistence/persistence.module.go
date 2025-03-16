package persistence

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewStockRepositoryImpl, fx.As(new(IStockRepository)))),
	fx.Provide(fx.Annotate(
		NewStockPersistanceAdapter,
		fx.As(new(port.IGetStockPort)), fx.As(new(port.IApplyStockUpdatePort)),
		fx.As(new(port.IApplyReservationEventPort)), fx.As(new(port.IGetReservationPort)),
	)),
	fx.Provide(fx.Annotate(NewCatalogRepositoryImpl, fx.As(new(ICatalogRepository)))),
	fx.Provide(fx.Annotate(NewCatalogPersistanceAdapter, fx.As(new(port.IApplyCatalogUpdatePort)), fx.As(new(port.IGetGoodPort)))),
)
