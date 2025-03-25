package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
)

type TransferRouter struct {
	transferController *TransferController
	broker             *broker.NatsMessageBroker
}

func NewTransferRouter(transferController *TransferController, broker *broker.NatsMessageBroker) *TransferRouter {
	return &TransferRouter{transferController, broker}
}

func (r *TransferRouter) Setup(ctx context.Context) error {
	// register request/reply handlers

	if err := r.broker.RegisterRequest(ctx, "transfer.create", "transfer", r.transferController.TransferCreateHandler); err != nil {
		return err
	}

	if err := r.broker.RegisterRequest(ctx, "transfer.get", "transfer", r.transferController.TransferGetHandler); err != nil {
		return err
	}

	if err := r.broker.RegisterRequest(ctx, "transfer.get.all", "transfer", r.transferController.TransferGetAllHandler); err != nil {
		return err
	}

	return nil
}
