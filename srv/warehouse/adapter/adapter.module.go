package adapter

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/listener"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/persistence"
	"github.com/alimitedgroup/MVP/srv/warehouse/adapter/sender"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controller.Module,
	persistence.Module,
	listener.Module,
	sender.Module,
)
