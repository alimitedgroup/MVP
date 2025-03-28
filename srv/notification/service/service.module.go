package service

import (
	serviceportin "github.com/alimitedgroup/MVP/srv/notification/service/portin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewNotificationService,
			fx.As(new(serviceportin.IAddQueryRuleUseCase)),
			fx.As(new(serviceportin.IAddStockUpdateUseCase)),
			fx.As(new(IService)),
		),
		NewRuleChecker,
	),

	fx.Invoke(func(rc *RuleChecker) {}),
)
