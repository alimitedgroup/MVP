package controller

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewOrderController),
	fx.Provide(NewHealthCheckController),
	fx.Provide(NewBrokerRoutes),
	fx.Provide(NewStockRouter),
	fx.Provide(NewHealthCheckRouter),
)
