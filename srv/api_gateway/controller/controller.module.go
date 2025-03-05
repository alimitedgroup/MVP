package controller

import (
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(AsController(NewHealthCheckController)),
	fx.Provide(AsController(NewLoginController)),
	fx.Invoke(fx.Annotate(RegisterRoutes, fx.ParamTags("", `group:"routes"`))),
)

func RegisterRoutes(http *lib.HTTPHandler, controllers []Controller) {
	for _, controller := range controllers {
		http.ApiGroup.Handle(controller.Method(), controller.Pattern(), controller.Handler())
	}
}

type Controller interface {
	Handler() gin.HandlerFunc
	Pattern() string
	Method() string
}

func AsController(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Controller)),
		fx.ResultTags(`group:"routes"`),
	)
}
