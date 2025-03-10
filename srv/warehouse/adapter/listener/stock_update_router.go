package listener

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
)

type StockUpdateRouter struct {
	stockListener *StockUpdateListener
	broker        *broker.NatsMessageBroker
	restore       broker.IRestoreStreamControl
	// StockUpdateConsumer jetstream.Consumer
}

func NewStockUpdateRouter(restoreFactory broker.IRestoreStreamControlFactory, stockUpdateListener *StockUpdateListener, n *broker.NatsMessageBroker) *StockUpdateRouter {
	return &StockUpdateRouter{stockUpdateListener, n, restoreFactory.Build()}
}

func (r *StockUpdateRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	err := r.broker.RegisterJsHandler(ctx, r.restore, stream.StockUpdateStreamConfig, r.stockListener.ListenStockUpdate)
	if err != nil {
		return err
	}
	// r.StockUpdateConsumer = consumer

	// wait restoring of the state before starting the server
	r.restore.Wait()

	// register request/reply handlers

	return nil
}
