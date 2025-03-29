package business

import (
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewBusiness,
			fx.As(new(portin.QueryRules)),
			fx.As(new(portin.StockUpdates)),
		),
		NewRuleChecker,
	),

	fx.Invoke(func(rc *RuleChecker) {}),
)
