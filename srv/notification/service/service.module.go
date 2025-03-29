package service

import (
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewNotificationService,
			fx.As(new(portin.IAddQueryRuleUseCase)),
			fx.As(new(portin.IAddStockUpdateUseCase)),
			fx.As(new(IService)),
		),
		NewRuleChecker,
	),

	fx.Invoke(func(rc *RuleChecker) {}),
)
