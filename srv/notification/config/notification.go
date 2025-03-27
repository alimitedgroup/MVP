package config

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/notification/controller"
	"github.com/alimitedgroup/MVP/srv/notification/notificationAdapter"
	"github.com/alimitedgroup/MVP/srv/notification/persistence"
	"github.com/alimitedgroup/MVP/srv/notification/service"
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

var Modules = fx.Options(
	lib.Module,
	controller.Module,
	notificationAdapter.Module,
	service.Module,
	persistence.Module,
)
