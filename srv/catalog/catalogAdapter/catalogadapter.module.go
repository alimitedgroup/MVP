package catalogAdapter

import (
	serviceportout "github.com/alimitedgroup/MVP/srv/catalog/service/portout"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewCatalogRepositoryAdapter,
			fx.As(new(serviceportout.IAddOrChangeGoodDataPort)),
			fx.As(new(serviceportout.ISetGoodQuantityPort)),
			fx.As(new(serviceportout.IGetGoodsQuantityPort)),
			fx.As(new(serviceportout.IGetGoodsInfoPort)),
			fx.As(new(serviceportout.IGetWarehousesInfoPort)))),
)
