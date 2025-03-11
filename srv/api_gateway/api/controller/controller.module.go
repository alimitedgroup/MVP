package controller

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthCheckController),
	fx.Provide(NewAPIRoutes),
	fx.Provide(NewHealthCheckRouter),
)
