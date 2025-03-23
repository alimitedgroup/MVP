package controller

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewAuthController),
	fx.Provide(NewAuthRouter),
	fx.Provide(NewControllerRouter),
)
