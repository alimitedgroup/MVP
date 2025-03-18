package persistence

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"go.uber.org/fx"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_persistence.go -package persistence github.com/alimitedgroup/MVP/srv/warehouse/adapter/persistence ICatalogRepository,IStockRepository,IIdempotentRepository

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewStockRepositoryImpl, fx.As(new(IStockRepository)))),
	fx.Provide(fx.Annotate(
		NewStockPersistanceAdapter,
		fx.As(new(port.IGetStockPort)), fx.As(new(port.IApplyStockUpdatePort)),
		fx.As(new(port.IApplyReservationEventPort)), fx.As(new(port.IGetReservationPort)),
	)),
	fx.Provide(fx.Annotate(NewCatalogRepositoryImpl, fx.As(new(ICatalogRepository)))),
	fx.Provide(fx.Annotate(NewCatalogPersistanceAdapter, fx.As(new(port.IApplyCatalogUpdatePort)), fx.As(new(port.IGetGoodPort)))),

	fx.Provide(fx.Annotate(NewIdempotentRepositoryImpl, fx.As(new(IIdempotentRepository)))),
	fx.Provide(fx.Annotate(NewIDempotentAdapter, fx.As(new(port.IIdempotentPort)))),
)
