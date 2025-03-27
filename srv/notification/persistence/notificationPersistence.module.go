package persistence

import (
	serviceportout "github.com/alimitedgroup/MVP/srv/notification/service/portout"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewRuleRepository,
			fx.As(new(serviceportout.IRuleRepository)),
		),
	),
)
