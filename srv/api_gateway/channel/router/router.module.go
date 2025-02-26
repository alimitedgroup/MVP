package router

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewStockUpdateRouter),
	fx.Provide(NewBrokerRoutes),
)

type BrokerRoute interface {
	Setup()
}

type BrokerRoutes []BrokerRoute

func NewBrokerRoutes(stockRouter *StockRouter) BrokerRoutes {
	return BrokerRoutes{
		stockRouter,
	}
}

func (r BrokerRoutes) Setup() {
	for _, v := range r {
		v.Setup()
	}
}
