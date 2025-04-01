package adapterout

import (
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewInfluxClient,
		fx.Annotate(NewNotificationAdapter,
			fx.As(new(portout.StockRepository)),
			fx.As(new(portout.StockEventPublisher)),
			fx.As(new(portout.RuleQueryRepository)),
		),
		fx.Annotate(NewRuleRepository, fx.As(new(portout.RuleRepository))),
	),
)
