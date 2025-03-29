package main

import (
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/notification/adapterout"
	"github.com/alimitedgroup/MVP/srv/notification/controller"
	"github.com/alimitedgroup/MVP/srv/notification/persistence"
	"github.com/alimitedgroup/MVP/srv/notification/service"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		lib.Module,
		controller.Module,
		adapterout.Module,
		service.Module,
		persistence.Module,
	)

	app.Run()
}
