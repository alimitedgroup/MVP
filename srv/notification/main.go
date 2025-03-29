package main

import (
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/notification/adapterin"
	"github.com/alimitedgroup/MVP/srv/notification/adapterout"
	"github.com/alimitedgroup/MVP/srv/notification/business"
	"github.com/alimitedgroup/MVP/srv/notification/persistence"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		lib.Module,
		adapterin.Module,
		adapterout.Module,
		business.Module,
		persistence.Module,
	)

	app.Run()
}
