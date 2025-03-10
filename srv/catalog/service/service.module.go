package service

import (
	serviceportin "github.com/alimitedgroup/MVP/srv/catalog/service/portin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewCatalogService,
		fx.As(new(serviceportin.IGetGoodsInfoUseCase)),
		fx.As(new(serviceportin.IGetGoodsQuantityUseCase)),
		fx.As(new(serviceportin.ISetMultipleGoodsQuantityUseCase)),
		fx.As(new(serviceportin.IUpdateGoodDataUseCase)),
		fx.As(new(serviceportin.IGetWarehousesUseCase)),
	)),
)
