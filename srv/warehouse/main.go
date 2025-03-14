package main

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"log"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/listener"
	"github.com/alimitedgroup/MVP/srv/warehouse/application"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"go.uber.org/fx"
)

type RunParams struct {
	fx.In

	BrokerRoutes   controller.BrokerRoutes
	ListenerRoutes listener.ListenerRoutes
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
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

var Modules = fx.Options(
	lib.Module,
	adapter.Module,
	application.Module,
)

func main() {
	ctx := context.Background()

	config := config.LoadConfig()

	opts := fx.Options(
		config,
		Modules,
		fx.Provide(broker.NewNatsConn),
	)

	app := fx.New(
		opts,
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
