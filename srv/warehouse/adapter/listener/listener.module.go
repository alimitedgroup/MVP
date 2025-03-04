package listener

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewStockUpdateListener),
	fx.Provide(NewStockUpdateRouter),
	fx.Provide(NewListenerRoutes),
)
