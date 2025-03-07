package application

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewApplyStockUpdateService, fx.As(new(port.ApplyStockUpdateUseCase)))),
	fx.Provide(fx.Annotate(NewManageStockService, fx.As(new(port.AddStockUseCase)), fx.As(new(port.RemoveStockUseCase)))),
	fx.Provide(fx.Annotate(NewApplyCatalogUpdateService, fx.As(new(port.ApplyCatalogUpdateUseCase)))),
)
