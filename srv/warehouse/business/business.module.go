package business

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"go.uber.org/fx"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_business.go -package=business github.com/alimitedgroup/MVP/srv/warehouse/business/port IApplyReservationEventPort,IApplyCatalogUpdatePort,IApplyStockUpdatePort,IStoreReservationEventPort,ICreateStockUpdatePort,IGetReservationPort,IGetGoodPort,IGetStockPort,IIdempotentPort,ITransactionPort

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewApplyStockUpdateService, fx.As(new(port.IApplyStockUpdateUseCase)))),
	fx.Provide(fx.Annotate(NewManageStockService, fx.As(new(port.IAddStockUseCase)), fx.As(new(port.IRemoveStockUseCase)))),
	fx.Provide(fx.Annotate(NewApplyCatalogUpdateService, fx.As(new(port.IApplyCatalogUpdateUseCase)))),
	fx.Provide(fx.Annotate(NewManageReservationService,
		fx.As(new(port.ICreateReservationUseCase)), fx.As(new(port.IApplyReservationUseCase)), fx.As(new(port.IConfirmOrderUseCase)), fx.As(new(port.IConfirmTransferUseCase)),
	)),
)
