package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
)

type StockRouter struct {
	stockController *StockController
	broker          *broker.NatsMessageBroker
	restore         *broker.RestoreStreamControl
}

func NewStockUpdateRouter(restore *broker.RestoreStreamControl, stockUpdateController *StockController, n *broker.NatsMessageBroker) *StockRouter {
	return &StockRouter{stockUpdateController, n, restore}
}

func (r *StockRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	err := r.broker.RegisterJsHandler(ctx, r.restore, stream.StockUpdateStreamConfig, r.stockController.UpdateHandler)
	if err != nil {
		return err
	}

	// wait restoring of the state before starting the server
	r.restore.Wait()

	return nil
}
