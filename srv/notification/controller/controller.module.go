package controller

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewNotificationController),
	fx.Provide(NewNotificationRouter),
	fx.Provide(NewControllerRouter),
)
