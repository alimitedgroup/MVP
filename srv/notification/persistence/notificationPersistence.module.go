package persistence

import (
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewRuleRepository,
			fx.As(new(portout.RuleRepository)),
		),
	),
)
