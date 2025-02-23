package router

import (
	"github.com/alimitedgroup/MVP/common/lib"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewHealthCheckRouter),
)

type Routes []lib.Route

func NewRoutes(healthCheckRoutes *HealthCheckRouter) Routes {
	return Routes{
		healthCheckRoutes,
	}
}

func (r Routes) Setup() {
	for _, v := range r {
		v.Setup()
	}
}
