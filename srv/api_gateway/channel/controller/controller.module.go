package controller

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewStockController),
	fx.Provide(NewStockUpdateRouter),
	fx.Provide(NewBrokerRoutes),
)
