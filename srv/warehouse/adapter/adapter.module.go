package adapter

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/listener"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/persistance"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/sender"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controller.Module,
	persistance.Module,
	listener.Module,
	sender.Module,
)
