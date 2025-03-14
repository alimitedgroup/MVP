package main

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"log"

	"github.com/alimitedgroup/MVP/srv/catalog/config"
	"go.uber.org/fx"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()
	app := fx.New(
		cfg,
		config.Modules,
		fx.Provide(broker.NewNatsConn),
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
