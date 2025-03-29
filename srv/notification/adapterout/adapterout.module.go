package adapterout

import (
	serviceportout2 "github.com/alimitedgroup/MVP/srv/notification/portout"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewInfluxClient,
		fx.Annotate(NewNotificationAdapter,
			fx.As(new(serviceportout2.IStockRepository)),
			fx.As(new(serviceportout2.IStockEventPublisher)),
			fx.As(new(serviceportout2.IRuleQueryRepository)),
		),
	),
)
