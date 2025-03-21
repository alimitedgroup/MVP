package listener

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
)

type ListenerRoutes []lib.BrokerRoute

func NewListenerRoutes(stockUpdateRouter *StockUpdateRouter, catalogRouter *CatalogRouter) ListenerRoutes {
	return ListenerRoutes{
		stockUpdateRouter,
		catalogRouter,
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
