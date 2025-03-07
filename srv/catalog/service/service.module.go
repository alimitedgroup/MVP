package service

import (
	service_portIn "github.com/alimitedgroup/MVP/srv/catalog/service/portIn"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewCatalogService,
		fx.As(new(service_portIn.IGetGoodsInfoUseCase)),
		fx.As(new(service_portIn.IGetGoodsQuantityUseCase)),
		fx.As(new(service_portIn.ISetMultipleGoodsQuantityUseCase)),
		fx.As(new(service_portIn.IUpdateGoodDataUseCase)),
		fx.As(new(service_portIn.IGetWarehousesUseCase)),
	)),
)
