package router

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/api_gateway/api/controller"
)

type HealthCheckRouter struct {
	http       *lib.HTTPHandler
	controller *controller.HealthCheckController
}

func NewHealthCheckRouter(http *lib.HTTPHandler, controller *controller.HealthCheckController) *HealthCheckRouter {
	return &HealthCheckRouter{http, controller}
}

func (p *HealthCheckRouter) Setup(ctx context.Context) {
	p.http.ApiGroup.GET("/ping", p.controller.Ping)
}
