package router

import (
	"github.com/alimitedgroup/MVP/common/lib"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewHealthCheckRouter),
)

type APIRoutes []lib.APIRoute

func NewRoutes(healthCheckRoutes *HealthCheckRouter) APIRoutes {
	return APIRoutes{
		healthCheckRoutes,
	}
}

func (r APIRoutes) Setup() {
	for _, v := range r {
		v.Setup()
	}
}
