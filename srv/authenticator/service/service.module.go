package service

import (
	serviceportin "github.com/alimitedgroup/MVP/srv/authenticator/service/portIn"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewAuthService,
			fx.As(new(serviceportin.IGetTokenUseCase)),
		),
	),
)
