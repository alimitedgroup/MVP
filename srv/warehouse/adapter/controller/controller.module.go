package controller

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewStockController),
	fx.Provide(NewHealthcheckController),
	fx.Provide(NewReservationController),
	fx.Provide(NewBrokerRoutes),
	fx.Provide(NewStockRouter),
	fx.Provide(NewHealthCheckRouter),
	fx.Provide(NewReservationRouter),
)
