package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewHTTPHandler),
	fx.Provide(AsController(NewHealthCheckController)),
	fx.Provide(AsController(NewLoginController)),
	fx.Provide(AsController(NewAuthHealthCheckController)),
	fx.Invoke(fx.Annotate(RegisterRoutes, fx.ParamTags("", `group:"routes"`))),
)

func RegisterRoutes(http *HTTPHandler, controllers []Controller) {
	for _, controller := range controllers {
		var group *gin.RouterGroup
		if controller.RequiresAuth() {
			group = http.AuthenticatedGroup
		} else {
			group = http.ApiGroup
		}
		group.Handle(controller.Method(), controller.Pattern(), controller.Handler())
	}
}

type Controller interface {
	Handler() gin.HandlerFunc
	Pattern() string
	Method() string
	RequiresAuth() bool
}

func AsController(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Controller)),
		fx.ResultTags(`group:"routes"`),
	)
}
