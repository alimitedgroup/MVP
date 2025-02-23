package router

import (
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/srv/api_gateway/api/controller"
)

type HealthCheckRouter struct {
	http       *lib.HTTPHandler
	controller *controller.HealthCheckController
}

func NewHealthCheckRouter(http *lib.HTTPHandler, controller *controller.HealthCheckController) lib.Route {
	return &HealthCheckRouter{http, controller}
}

func (p *HealthCheckRouter) Setup() {
	p.http.ApiGroup.GET("/ping", p.controller.Ping)
}
