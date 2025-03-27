package notificationAdapter

import (
	serviceportout "github.com/alimitedgroup/MVP/srv/notification/service/portout"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewStockAdapter,
			fx.As(new(serviceportout.IStockRepository)),
			fx.As(new(serviceportout.IStockEventPublisher)),
			fx.As(new(serviceportout.IRuleQueryRepository)),
		),
	),
)
