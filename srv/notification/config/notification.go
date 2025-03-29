package config

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/notification/controller"
	"go.uber.org/fx"
)

func Run(ctx context.Context, nr *controller.ControllerRouter) error {
	err := nr.Setup(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func RunLifeCycle(lc fx.Lifecycle, nr *controller.ControllerRouter) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return Run(ctx, nr)
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
