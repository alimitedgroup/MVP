package business

import (
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"go.uber.org/fx"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_business.go -package=business github.com/alimitedgroup/MVP/srv/order/business/port ISetCompleteTransferPort,ISetCompletedWarehouseOrderPort,ISendContactWarehousePort,IApplyOrderUpdatePort,IApplyStockUpdatePort,IApplyTransferUpdatePort,IGetOrderPort,IGetStockPort,IGetTransferPort,ISendOrderUpdatePort,IRequestReservationPort,ISendTransferUpdatePort,ICalculateAvailabilityUseCase

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewManageOrderService,
		fx.As(new(port.ICreateOrderUseCase)), fx.As(new(port.IGetOrderUseCase)), fx.As(new(port.IContactWarehousesUseCase)),
		fx.As(new(port.ICreateTransferUseCase)), fx.As(new(port.IGetTransferUseCase)),
	)),
	fx.Provide(fx.Annotate(NewApplyOrderUpdateService, fx.As(new(port.IApplyOrderUpdateUseCase)), fx.As(new(port.IApplyTransferUpdateUseCase)))),
	fx.Provide(fx.Annotate(NewApplyStockUpdateService, fx.As(new(port.IApplyStockUpdateUseCase)))),
	fx.Provide(fx.Annotate(NewSimpleCalculateAvailabilityService, fx.As(new(port.ICalculateAvailabilityUseCase)))),
)
