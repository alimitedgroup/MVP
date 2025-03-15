package main

import (
	"context"
	"log"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/order/adapter"
	"github.com/alimitedgroup/MVP/srv/order/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/order/adapter/listener"
	"github.com/alimitedgroup/MVP/srv/order/business"
	"github.com/alimitedgroup/MVP/srv/order/config"
	"go.uber.org/fx"
)

type RunParams struct {
	fx.In

	BrokerRoutes   controller.ControllerRoutes
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
	business.Module,
)

func main() {
	ctx := context.Background()

	config := config.LoadConfig()

	opts := fx.Options(
		config,
		Modules,
	)

	app := fx.New(
		opts,
		fx.Provide(broker.NewNatsConn),
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
