package persistence

import (
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"go.uber.org/fx"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_persistence.go -package=persistence github.com/alimitedgroup/MVP/srv/order/adapter/persistence IStockRepository,ITransferRepository,IOrderRepository

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewStockRepositoryImpl, fx.As(new(IStockRepository)))),
	fx.Provide(fx.Annotate(NewTransferRepositoryImpl, fx.As(new(ITransferRepository)))),
	fx.Provide(fx.Annotate(NewOrderPersistanceAdapter,
		fx.As(new(port.IGetOrderPort)), fx.As(new(port.IApplyOrderUpdatePort)), fx.As(new(port.ISetCompletedWarehouseOrderPort)),
	)),
	fx.Provide(fx.Annotate(NewTransferPersistanceAdapter,
		fx.As(new(port.IGetTransferPort)), fx.As(new(port.IApplyTransferUpdatePort)), fx.As(new(port.ISetCompleteTransferPort)),
	)),
	fx.Provide(fx.Annotate(NewOrderRepositoryImpl, fx.As(new(IOrderRepository)))),
	fx.Provide(fx.Annotate(NewStockPersistanceAdapter, fx.As(new(port.IGetStockPort)), fx.As(new(port.IApplyStockUpdatePort)))),
)
