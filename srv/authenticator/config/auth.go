package config

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/authenticator/adapter"
	"github.com/alimitedgroup/MVP/srv/authenticator/controller"
	"github.com/alimitedgroup/MVP/srv/authenticator/persistence"
	"github.com/alimitedgroup/MVP/srv/authenticator/publisher"
	"github.com/alimitedgroup/MVP/srv/authenticator/service"
	serviceauthenticator "github.com/alimitedgroup/MVP/srv/authenticator/service/strategy"
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
	adapter.Module,
	controller.Module,
	persistence.Module,
	publisher.Module,
	service.Module,
	serviceauthenticator.Module,
)

var ModulesTest = fx.Options(
	lib.ModuleTest,
	adapter.Module,
	controller.Module,
	persistence.Module,
	publisher.Module,
	service.Module,
	serviceauthenticator.Module,
)
