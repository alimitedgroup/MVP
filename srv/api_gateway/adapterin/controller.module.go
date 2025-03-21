package adapterin

import (
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"adapterin",
	fx.Decorate(observability.WrapLogger("adapterin")),
	fx.Provide(NewHTTPHandler),
	fx.Provide(AsController(NewHealthCheckController)),
	fx.Provide(AsController(NewLoginController)),
	fx.Provide(AsController(NewAuthHealthCheckController)),
	fx.Provide(AsController(NewListWarehousesController)),
	fx.Provide(AsController(NewGetGoodsController)),
	fx.Invoke(fx.Annotate(RegisterRoutes, fx.ParamTags("", `group:"routes"`))),
)

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
