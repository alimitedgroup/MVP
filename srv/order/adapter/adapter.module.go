package adapter

import (
	"github.com/alimitedgroup/MVP/srv/order/adapter/controller"
	"github.com/alimitedgroup/MVP/srv/order/adapter/listener"
	"github.com/alimitedgroup/MVP/srv/order/adapter/persistence"
	"github.com/alimitedgroup/MVP/srv/order/adapter/sender"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controller.Module,
	persistence.Module,
	listener.Module,
	sender.Module,
)
