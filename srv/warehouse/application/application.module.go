package application

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewApplyStockUpdateService, fx.As(new(port.IApplyStockUpdateUseCase)))),
	fx.Provide(fx.Annotate(NewManageStockService, fx.As(new(port.IAddStockUseCase)), fx.As(new(port.IRemoveStockUseCase)))),
	fx.Provide(fx.Annotate(NewApplyCatalogUpdateService, fx.As(new(port.IApplyCatalogUpdateUseCase)))),
)
