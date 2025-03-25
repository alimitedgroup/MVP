package listener

import (
	"context"
)

type ListenerRoutes struct {
	StockUpdateRouter *StockUpdateRouter
	OrderUpdateRouter *OrderRouter
}

func NewListenerRoutes(stockUpdateRouter *StockUpdateRouter, orderUpdateRouter *OrderRouter) *ListenerRoutes {
	return &ListenerRoutes{
		stockUpdateRouter,
		orderUpdateRouter,
	}
}

func (r ListenerRoutes) Setup(ctx context.Context) error {
	if err := r.OrderUpdateRouter.Setup(ctx); err != nil {
		return err
	}

	if err := r.StockUpdateRouter.Setup(ctx); err != nil {
		return err
	}

	return nil
}
