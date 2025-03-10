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

func NewStockRouter(config *config.WarehouseConfig, stockUpdateController *StockController, broker *broker.NatsMessageBroker) *StockRouter {
	return &StockRouter{config, stockUpdateController, broker}
}

func (r *StockRouter) Setup(ctx context.Context) error {
	// register request/reply handlers
	var err error

	err = r.broker.RegisterRequest(ctx, broker.Subject(fmt.Sprintf("warehouse.%s.stock.add", r.config.ID)), broker.NoQueue, r.stockController.AddStockHandler)
	if err != nil {
		return err
	}

	err = r.broker.RegisterRequest(ctx, broker.Subject(fmt.Sprintf("warehouse.%s.stock.remove", r.config.ID)), broker.NoQueue, r.stockController.RemoveStockHandler)
	if err != nil {
		return err
	}

	return nil
}
