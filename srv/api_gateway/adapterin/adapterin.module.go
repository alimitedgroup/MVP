package adapterin

import (
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var ModuleTest = fx.Module(
	"adapterin",
	fx.Decorate(observability.WrapLogger("adapterin")),
	fx.Provide(NewHTTPHandler),
	fx.Provide(AsController(NewHealthCheckController)),
	fx.Provide(AsController(NewLoginController)),
	fx.Provide(AsController(NewAuthHealthCheckController)),
	fx.Provide(AsController(NewListWarehousesController)),
	fx.Provide(AsController(NewGetGoodsController)),
	fx.Provide(AsController(NewCreateGoodController)),
	fx.Provide(AsController(NewUpdateGoodController)),
	fx.Provide(AsController(NewGetOrdersController)),
	fx.Provide(AsController(NewCreateOrderController)),
	fx.Provide(AsController(NewGetTransfersController)),
	fx.Provide(AsController(NewCreateTransferController)),
	fx.Provide(AsController(NewAddStockController)),
	fx.Provide(AsController(NewRemoveStockController)),
	fx.Invoke(fx.Annotate(RegisterRoutes, fx.ParamTags("", `group:"routes"`))),
)

var Module = fx.Module(
	"adapterin",
	fx.Decorate(observability.WrapLogger("adapterin")),
	fx.Provide(ConfigFromEnv, NewListener, NewHTTPHandler),
	fx.Provide(AsController(NewHealthCheckController)),
	fx.Provide(AsController(NewLoginController)),
	fx.Provide(AsController(NewAuthHealthCheckController)),
	fx.Provide(AsController(NewListWarehousesController)),
	fx.Provide(AsController(NewGetGoodsController)),
	fx.Provide(AsController(NewCreateGoodController)),
	fx.Provide(AsController(NewUpdateGoodController)),
	fx.Provide(AsController(NewGetOrdersController)),
	fx.Provide(AsController(NewCreateOrderController)),
	fx.Provide(AsController(NewGetTransfersController)),
	fx.Provide(AsController(NewCreateTransferController)),
	fx.Provide(AsController(NewAddStockController)),
	fx.Provide(AsController(NewRemoveStockController)),
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
