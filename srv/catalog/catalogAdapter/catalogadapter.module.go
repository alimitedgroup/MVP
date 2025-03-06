package catalogAdapter

import (
	service_portOut "github.com/alimitedgroup/MVP/srv/catalog/service/portOut"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewCatalogRepositoryAdapter,
			fx.As(new(service_portOut.IAddOrChangeGoodDataPort)),
			fx.As(new(service_portOut.ISetGoodQuantityPort)),
			fx.As(new(service_portOut.IGetGoodsQuantityPort)),
			fx.As(new(service_portOut.IGetGoodsInfoPort)))),
)
