package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
)

type HealthCheckRouter struct {
	http       *lib.HTTPHandler
	controller *HealthCheckController
}

func NewHealthCheckRouter(http *lib.HTTPHandler, controller *HealthCheckController) *HealthCheckRouter {
	return &HealthCheckRouter{http, controller}
}

func (p *HealthCheckRouter) Setup(ctx context.Context) {
	p.http.ApiGroup.GET("/ping", p.controller.Ping)
}
