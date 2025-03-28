package main

import (
	"context"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"log"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/listener"
	"github.com/alimitedgroup/MVP/srv/warehouse/business"
	"go.uber.org/fx"
)

type RunParams struct {
	fx.In

	BrokerRoutes   *controller.BrokerRoutes
	ListenerRoutes *listener.ListenerRoutes
}

func Run(ctx context.Context, p RunParams) error {
	var err error

	err = p.ListenerRoutes.Setup(ctx)
	if err != nil {
		return err
	}

	err = p.BrokerRoutes.Setup(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func RunLifeCycle(lc fx.Lifecycle, p RunParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := Run(ctx, p)
			return err
		},
	})
}

var Modules = fx.Options(
	lib.Module,
	adapter.Module,
	business.Module,
)

func main() {
	ctx := context.Background()

	app := fx.New(
		Modules,
		fx.Provide(config.ConfigFromEnv),
		fx.Invoke(RunLifeCycle),
	)

	err := app.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

}
