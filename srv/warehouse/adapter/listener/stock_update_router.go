package listener

import (
	"context"
	"fmt"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
)

type StockUpdateRouter struct {
	stockListener *StockUpdateListener
	broker        *broker.NatsMessageBroker
	restore       broker.IRestoreStreamControl
	cfg           *config.WarehouseConfig
}

func NewStockUpdateRouter(restoreFactory broker.IRestoreStreamControlFactory, stockUpdateListener *StockUpdateListener, broker *broker.NatsMessageBroker, cfg *config.WarehouseConfig) *StockUpdateRouter {
	return &StockUpdateRouter{stockUpdateListener, broker, restoreFactory.Build(), cfg}
}

func (r *StockUpdateRouter) Setup(ctx context.Context) error {
	// register stream message handlers
	err := r.broker.RegisterJsHandler(ctx, r.restore, stream.StockUpdateStreamConfig, r.stockListener.ListenStockUpdate, broker.WithSubjectFilter(fmt.Sprintf("stock.update.%s", r.cfg.ID)))
	if err != nil {
		return err
	}

	// wait restoring of the state before starting the server
	r.restore.Wait()

	return nil
}
