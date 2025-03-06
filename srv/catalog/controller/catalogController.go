package controller

import (
	"context"

	service_portIn "github.com/alimitedgroup/MVP/srv/catalog/service/portIn"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type catalogController struct {
	getGoodsInfoUseCase             service_portIn.IGetGoodsInfoUseCase
	getGoodsQuantityUseCase         service_portIn.IGetGoodsQuantityUseCase
	setMultipleGoodsQuantityUseCase service_portIn.ISetMultipleGoodsQuantityUseCase
	updateGoodDataUseCase           service_portIn.IUpdateGoodDataUseCase
}

func NewCatalogController(getGoodsInfoUseCase service_portIn.IGetGoodsInfoUseCase, getGoodsQuantityUseCase service_portIn.IGetGoodsQuantityUseCase, setMultipleGoodsQuantityUseCase service_portIn.ISetMultipleGoodsQuantityUseCase, updateGoodDataUseCase service_portIn.IUpdateGoodDataUseCase) *catalogController {
	return &catalogController{getGoodsInfoUseCase: getGoodsInfoUseCase, getGoodsQuantityUseCase: getGoodsQuantityUseCase, setMultipleGoodsQuantityUseCase: setMultipleGoodsQuantityUseCase, updateGoodDataUseCase: updateGoodDataUseCase}
}

func (cc *catalogController) getGoodRequest(ctx context.Context, msg *nats.Msg) {

}

func (cc *catalogController) getWarehouseRequest(ctx context.Context, msg *nats.Msg) {
}

func (cc *catalogController) setGoodDataRequest(ctx context.Context, msg jetstream.Msg) {
}

func (cc *catalogController) setGoodQuantityRequest(ctx context.Context, msg jetstream.Msg) {
}
