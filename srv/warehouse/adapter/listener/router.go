package listener

import (
	"context"
)

type ListenerRoutes struct {
	StockUpdateRouter *StockUpdateRouter
	CatalogRouter     *CatalogRouter
	ReservationRouter *ReservationEventRouter
	OrderUpdateRouter *OrderUpdateRouter
}

func NewListenerRoutes(
	stockUpdateRouter *StockUpdateRouter, catalogRouter *CatalogRouter,
	reservationRouter *ReservationEventRouter, orderUpdateRouter *OrderUpdateRouter,
) *ListenerRoutes {
	return &ListenerRoutes{
		stockUpdateRouter,
		catalogRouter,
		reservationRouter,
		orderUpdateRouter,
	}
}

// NOTE: must setup() in this order
func (r ListenerRoutes) Setup(ctx context.Context) error {
	if err := r.CatalogRouter.Setup(ctx); err != nil {
		return err
	}
	if err := r.StockUpdateRouter.Setup(ctx); err != nil {
		return err
	}
	if err := r.ReservationRouter.Setup(ctx); err != nil {
		return err
	}
	if err := r.OrderUpdateRouter.Setup(ctx); err != nil {
		return err
	}

	return nil
}
