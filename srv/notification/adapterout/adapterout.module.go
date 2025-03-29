package adapterout

import (
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewInfluxClient,
		fx.Annotate(NewNotificationAdapter,
			fx.As(new(portout.IStockRepository)),
			fx.As(new(portout.IStockEventPublisher)),
			fx.As(new(portout.IRuleQueryRepository)),
		),
		fx.Annotate(NewRuleRepository, fx.As(new(portout.RuleRepository))),
	),
)
