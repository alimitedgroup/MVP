package catalogAdapter

import (
	"github.com/alimitedgroup/MVP/srv/catalog/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewCatalogRepositoryAdapter,
			fx.As(new(service.IAddOrChangeGoodDataPort)),
			fx.As(new(service.ISetGoodQuantityPort)),
			fx.As(new(service.IGetGoodsQuantityPort)),
			fx.As(new(service.IGetGoodsInfoPort)))),
)
