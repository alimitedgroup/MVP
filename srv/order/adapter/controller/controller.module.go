package controller

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthCheckController),
	fx.Provide(NewHealthCheckRouter),
	fx.Provide(NewOrderController),
	fx.Provide(NewOrderRouter),
	fx.Provide(NewTransferController),
	fx.Provide(NewTransferRouter),
	fx.Provide(NewBrokerRoutes),
)
