package persistence

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewAuthRepo,
		fx.As(new(IAuthPersistance)),
	)),
)
