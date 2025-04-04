package main

import (
	"context"
	"log"

	"github.com/alimitedgroup/MVP/srv/catalog/config"
	"go.uber.org/fx"
)

func main() {
	ctx := context.Background()
	app := fx.New(
		config.Modules,
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
