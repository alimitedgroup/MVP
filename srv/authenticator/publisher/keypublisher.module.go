package publisher

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewPublisher,
		fx.As(new(IAuthPublisher)),
	)),
)
