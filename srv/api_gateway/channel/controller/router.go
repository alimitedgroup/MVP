package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
)

type BrokerRoutes []lib.BrokerRoute

func NewBrokerRoutes(stockRouter *StockRouter) BrokerRoutes {
	return BrokerRoutes{
		stockRouter,
	}
}

func (r BrokerRoutes) Setup(ctx context.Context) error {
	for _, v := range r {
		err := v.Setup(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
