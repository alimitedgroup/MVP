package controller

import (
	"log"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go"
)

type StockController struct {
	n *broker.NatsMessageBroker
}

func NewStockController(n *broker.NatsMessageBroker) *StockController {
	return &StockController{n}
}

func (c *StockController) UpdateHandler(msg *nats.Msg) {
	log.Printf("Received a message: %s\n", string(msg.Data))
}
