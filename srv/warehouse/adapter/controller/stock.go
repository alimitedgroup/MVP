package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go"
)

type StockController struct {
	broker *broker.NatsMessageBroker
}

func NewStockController(n *broker.NatsMessageBroker) *StockController {
	return &StockController{n}
}

func (c *StockController) AddStockHandler(ctx context.Context, msg *nats.Msg) error {

	return nil
}

func (c *StockController) RemoveStockHandler(ctx context.Context, msg *nats.Msg) error {
	return nil
}
