package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
)

type OrderRouter struct {
	orderController *OrderController
	broker          *broker.NatsMessageBroker
}

func NewStockRouter(orderController *OrderController, broker *broker.NatsMessageBroker) *OrderRouter {
	return &OrderRouter{orderController, broker}
}

func (r *OrderRouter) Setup(ctx context.Context) error {
	// register request/reply handlers
	err := r.broker.RegisterRequest(ctx, "order.create", broker.NoQueue, r.orderController.OrderCreateHandler)
	if err != nil {
		return err
	}

	return nil
}
