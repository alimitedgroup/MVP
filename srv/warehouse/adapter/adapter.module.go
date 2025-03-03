package adapter

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/listener"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/persistance"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controller.Module,
	persistance.Module,
	listener.Module,
)
