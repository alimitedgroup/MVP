package main

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/api_gateway/api"
	"github.com/alimitedgroup/MVP/srv/api_gateway/api/router"
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel"
	"go.uber.org/fx"
)

func Run(h *lib.HTTPHandler, routes router.APIRoutes) {
	routes.Setup()

	_ = h.Engine.Run(":8080")
}

var Modules = fx.Options(
	lib.Module,
	api.Module,
	channel.Module,
)

func main() {
	_ = context.Background()

	opts := fx.Options(Modules)
	app := fx.New(
		opts,
		fx.Invoke(Run),
	)

	app.Run()
}
