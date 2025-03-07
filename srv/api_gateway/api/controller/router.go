package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
)

type APIRoutes []lib.APIRoute

func NewAPIRoutes(healthCheckRoutes *HealthCheckRouter) APIRoutes {
	return APIRoutes{
		healthCheckRoutes,
	}
}

func (r APIRoutes) Setup(ctx context.Context) {
	for _, v := range r {
		v.Setup(ctx)
	}
}
