package serviceauthenticator

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewSimpleAuthAlg,
		fx.As(new(IAuthenticateUserStrategy)),
	)),
)
