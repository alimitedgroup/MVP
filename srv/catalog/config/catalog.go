package config

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogAdapter"
	"github.com/alimitedgroup/MVP/srv/catalog/controller"
	goodRepository "github.com/alimitedgroup/MVP/srv/catalog/persistence"
	"github.com/alimitedgroup/MVP/srv/catalog/service"
	"go.uber.org/fx"
)

func Run(ctx context.Context, cr *controller.ControllerRouter) error {
	//var err error
	err := cr.Setup(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func RunLifeCycle(lc fx.Lifecycle, cr *controller.ControllerRouter) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := Run(ctx, cr)
			return err
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

var Modules = fx.Options(
	lib.Module,
	controller.Module,
	goodRepository.Module,
	catalogAdapter.Module,
	service.Module,
)
