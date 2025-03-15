package listener

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
)

type ListenerRoutes []lib.BrokerRoute

func NewListenerRoutes(stockUpdateRouter *StockUpdateRouter, orderUpdateRouter *OrderUpdateRouter) ListenerRoutes {
	return ListenerRoutes{
		stockUpdateRouter,
		orderUpdateRouter,
	}
}

func (r ListenerRoutes) Setup(ctx context.Context) error {
	for _, v := range r {
		err := v.Setup(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
