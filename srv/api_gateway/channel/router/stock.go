package router

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel/constants"
	"github.com/alimitedgroup/MVP/srv/api_gateway/channel/controller"
)

type StockRouter struct {
	stockController *controller.StockController
	broker          *broker.NatsMessageBroker
	restore         *broker.RestoreStreamControl
}

func NewStockUpdateRouter(restore *broker.RestoreStreamControl, stockUpdateController *controller.StockController, n *broker.NatsMessageBroker) *StockRouter {
	return &StockRouter{stockUpdateController, n, restore}
}

func (r *StockRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	err := r.broker.RegisterJsHandler(ctx, r.restore, constants.StockUpdatesStreamConfig, r.stockController.UpdateHandler)
	if err != nil {
		return err
	}

	// wait restoring of the state before starting the server
	r.restore.Wait()

	// register request/reply handlers

	return nil
}
