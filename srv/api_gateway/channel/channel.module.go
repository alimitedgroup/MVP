package channel

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel/controller"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controller.Module,
)
