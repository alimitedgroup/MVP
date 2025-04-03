package controller

import (
	"context"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
)

type catalogRouter struct {
	mb             *broker.NatsMessageBroker
	controller     *catalogController
	goodController *CatalogGoodInfoController
	qtController   *CatalogGlobalQuantityController
	rsc            *broker.RestoreStreamControl
}

func NewCatalogRouter(mb *broker.NatsMessageBroker, cc *catalogController, gc *CatalogGoodInfoController, qt *CatalogGlobalQuantityController, rsc *broker.RestoreStreamControl) *catalogRouter {
	return &catalogRouter{mb, cc, gc, qt, rsc}
}

func (cr *catalogRouter) Setup(ctx context.Context) error {
	err := cr.mb.RegisterJsHandler(ctx, cr.rsc, stream.StockUpdateStreamConfig, cr.controller.SetGoodQuantityRequest) //SetMultipleGoodsQuantity
	if err != nil {
		return nil
	}
	err = cr.mb.RegisterJsHandler(ctx, cr.rsc, stream.AddOrChangeGoodDataStream, cr.goodController.SetGoodDataRequest) //AddOrChangeGoodData
	if err != nil {
		return nil
	}
	cr.rsc.Wait()
	err = cr.mb.RegisterRequest(ctx, "catalog.getGoods", "catalog", cr.goodController.GetGoodsRequest) //GetGoodsInfo
	if err != nil {
		return nil
	}
	err = cr.mb.RegisterRequest(ctx, "catalog.getWarehouses", "catalog", cr.controller.GetWarehouseRequest) //GetWarehouses
	if err != nil {
		return nil
	}
	err = cr.mb.RegisterRequest(ctx, "catalog.getGoodsGlobalQuantity", "catalog", cr.qtController.GetGoodsGlobalQuantityRequest) //GetGoodsQuantity
	if err != nil {
		return nil
	}
	return nil
}
