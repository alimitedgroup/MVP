package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
)

type catalogRouter struct {
	mb         *broker.NatsMessageBroker
	controller *catalogController
	rsc        *broker.RestoreStreamControl
}

func NewCatalogRouter(mb *broker.NatsMessageBroker, cc *catalogController, rsc *broker.RestoreStreamControl) *catalogRouter {
	return &catalogRouter{mb, cc, rsc}
}

func (cr *catalogRouter) Setup(ctx context.Context) error {
	cr.mb.RegisterJsHandler(ctx, cr.rsc, stream.StockUpdateStreamConfig, cr.controller.setGoodQuantityRequest)
	cr.mb.RegisterJsHandler(ctx, cr.rsc, stream.StockUpdateStreamConfig, cr.controller.setGoodDataRequest) //sistemare con handler giusto
	cr.rsc.Wait()
	cr.mb.RegisterRequest(ctx, "catalog.getGoods", "catalog", cr.controller.getGoodRequest)
	cr.mb.RegisterRequest(ctx, "catalog.getWarehouse", "catalog", cr.controller.getWarehouseRequest)
	return nil
}
