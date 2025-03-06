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
	err := cr.mb.RegisterJsHandler(ctx, cr.rsc, stream.StockUpdateStreamConfig, cr.controller.setGoodQuantityRequest) //SetMultipleGoodsQuantity
	if err != nil {
		return nil
	}
	err = cr.mb.RegisterJsHandler(ctx, cr.rsc, stream.AddOrChangeGoodDataStream, cr.controller.setGoodDataRequest) //AddOrChangeGoodData
	if err != nil {
		return nil
	}
	cr.rsc.Wait()
	err = cr.mb.RegisterRequest(ctx, "catalog.getGoods", "catalog", cr.controller.getGoodsRequest) //GetGoodsInfo
	if err != nil {
		return nil
	}
	err = cr.mb.RegisterRequest(ctx, "catalog.getWarehouse", "catalog", cr.controller.getWarehouseRequest) //GetWarehouses
	if err != nil {
		return nil
	}
	err = cr.mb.RegisterRequest(ctx, "catalog.getGoodsGlobalQuantity", "catalog", cr.controller.getGoodsGlobalQuantityRequest) //GetGoodsQuantity
	if err != nil {
		return nil
	}
	return nil
}
