package router

import (
	"github.com/alimitedgroup/MVP/common/lib"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewAPIRoutes),
	fx.Provide(NewHealthCheckRouter),
)

type BrokerRoute interface {
	Setup()
}

type APIRoutes []lib.APIRoute

func NewAPIRoutes(healthCheckRoutes *HealthCheckRouter) APIRoutes {
	return APIRoutes{
		healthCheckRoutes,
	}
}

func (r APIRoutes) Setup() {
	for _, v := range r {
		v.Setup()
	}
}
