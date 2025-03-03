package api

import (
	"github.com/alimitedgroup/MVP/srv/api_gateway/api/controller"
	"github.com/alimitedgroup/MVP/srv/api_gateway/api/router"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controller.Module,
	router.Module,
)
