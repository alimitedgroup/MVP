package business

import (
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewManageStockService, fx.As(new(port.ICreateOrderUseCase)))),
	fx.Provide(fx.Annotate(NewApplyStockUpdateService, fx.As(new(port.IApplyStockUpdateUseCase)))),
	fx.Provide(fx.Annotate(NewSimpleCalculateAvailabilityService, fx.As(new(port.ICalculateAvailabilityUseCase)))),
)
