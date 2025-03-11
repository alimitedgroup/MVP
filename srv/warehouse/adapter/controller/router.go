package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib"
)

type BrokerRoutes []lib.BrokerRoute

func NewBrokerRoutes(stockRouter *StockRouter, reservationRouter *ReservationRouter, healthCheckRouter *HealthCheckRouter) BrokerRoutes {
	return BrokerRoutes{
		stockRouter,
		reservationRouter,
		healthCheckRouter,
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
