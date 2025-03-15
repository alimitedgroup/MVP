package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

type TransferController struct {
	broker *broker.NatsMessageBroker
}

type TransferControllerParams struct {
	fx.In

	Broker *broker.NatsMessageBroker
}

func NewTransferController(p TransferControllerParams) *TransferController {
	return &TransferController{p.Broker}
}

func (c *TransferController) TransferCreateHandler(ctx context.Context, msg *nats.Msg) error {
	// implementation
	return nil
}

func (c *TransferController) TransferGetHandler(ctx context.Context, msg *nats.Msg) error {
	// implementation
	return nil
}

func (c *TransferController) TransferGetAllHandler(ctx context.Context, msg *nats.Msg) error {
	// implementation
	return nil
}
