package controller

import (
	"context"
	"fmt"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
)

type StockRouter struct {
	config          *config.WarehouseConfig
	stockController *StockController
	broker          *broker.NatsMessageBroker
}

func NewStockUpdateRouter(config *config.WarehouseConfig, stockUpdateController *StockController, broker *broker.NatsMessageBroker) *StockRouter {
	return &StockRouter{config, stockUpdateController, broker}
}

func (r *StockRouter) Setup(ctx context.Context) error {
	// register request/reply handlers
	r.broker.RegisterRequest(ctx, broker.Subject(fmt.Sprintf("warehouse.stock.add.%s", r.config.ID)), broker.NoQueue, r.stockController.AddStockHandler)
	r.broker.RegisterRequest(ctx, broker.Subject(fmt.Sprintf("warehouse.stock.remove.%s", r.config.ID)), broker.NoQueue, r.stockController.RemoveStockHandler)

	return nil
}
