package main

import (
	"context"
	"log"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/catalog/config"
	"github.com/alimitedgroup/MVP/srv/catalog/controller"
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
	//metti tutti altri module
)

func main() {
	ctx := context.Background()
	config := config.LoadConfig()
	app := fx.New(
		Modules,
		config,
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
