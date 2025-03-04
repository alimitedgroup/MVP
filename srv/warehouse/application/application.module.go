package application

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewUpdateStockService,
		fx.As(new(port.UpdateStockUseCase))),
	),
)
