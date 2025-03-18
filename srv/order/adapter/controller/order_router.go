package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
)

type OrderRouter struct {
	orderController *OrderController
	broker          *broker.NatsMessageBroker
}

func NewOrderRouter(orderController *OrderController, broker *broker.NatsMessageBroker) *OrderRouter {
	return &OrderRouter{orderController, broker}
}

func (r *OrderRouter) Setup(ctx context.Context) error {
	// register request/reply handlers

	if err := r.broker.RegisterRequest(ctx, "order.create", "order", r.orderController.OrderCreateHandler); err != nil {
		return err
	}

	if err := r.broker.RegisterRequest(ctx, "order.get", "order", r.orderController.OrderGetHandler); err != nil {
		return err
	}

	if err := r.broker.RegisterRequest(ctx, "order.get.all", "order", r.orderController.OrderGetAllHandler); err != nil {
		return err
	}

	return nil
}
