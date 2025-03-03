package channel

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel/controller"
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel/router"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controller.Module,
	router.Module,
)
