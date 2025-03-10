package api

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/api/controller"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controller.Module,
)
