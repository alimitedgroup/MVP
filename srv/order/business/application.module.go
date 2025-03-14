package business

import (
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewManageStockService, fx.As(new(port.ICreateOrderUseCase)), fx.As(new(port.IGetOrderUseCase)))),
	fx.Provide(fx.Annotate(NewApplyOrderUpdateService, fx.As(new(port.IApplyOrderUpdateUseCase)))),
	fx.Provide(fx.Annotate(NewApplyStockUpdateService, fx.As(new(port.IApplyStockUpdateUseCase)))),
	fx.Provide(fx.Annotate(NewSimpleCalculateAvailabilityService, fx.As(new(port.ICalculateAvailabilityUseCase)))),
)
