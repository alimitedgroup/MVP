package main

import (
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/api_gateway/adapterin"
	"github.com/alimitedgroup/MVP/srv/api_gateway/adapterout"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		lib.Module,
		business.Module,
		adapterout.Module,
		adapterin.Module,
	)

	app.Run()
}
