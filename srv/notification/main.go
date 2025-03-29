package main

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/notification/config"
	"github.com/alimitedgroup/MVP/srv/notification/controller"
	"github.com/alimitedgroup/MVP/srv/notification/notificationAdapter"
	"github.com/alimitedgroup/MVP/srv/notification/persistence"
	"github.com/alimitedgroup/MVP/srv/notification/service"
	"log"

	"go.uber.org/fx"
)

func main() {
	ctx := context.Background()

	app := fx.New(
		lib.Module,
		controller.Module,
		notificationAdapter.Module,
		service.Module,
		persistence.Module,
		fx.Invoke(config.RunLifeCycle),
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
