package router

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewStockUpdateRouter),
	fx.Provide(NewBrokerRoutes),
)

type BrokerRoutes []lib.BrokerRoute

func NewBrokerRoutes(stockRouter *StockRouter) BrokerRoutes {
	return BrokerRoutes{
		stockRouter,
	}
}

func (r BrokerRoutes) Setup(ctx context.Context) {
	for _, v := range r {
		v.Setup(ctx)
	}
}
