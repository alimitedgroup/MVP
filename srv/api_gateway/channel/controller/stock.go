package controller

import (
	"context"
	"log"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go/jetstream"
)

type StockController struct {
	n *broker.NatsMessageBroker
}

func NewStockController(n *broker.NatsMessageBroker) *StockController {
	return &StockController{n}
}

func (c *StockController) UpdateHandler(ctx context.Context, msg jetstream.Msg) error {
	log.Printf("Received a message: %s\n", string(msg.Data()))
	return nil
}
