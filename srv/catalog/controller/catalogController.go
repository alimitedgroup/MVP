package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/catalog/service/portin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

type catalogController struct {
	getGoodsInfoUseCase             serviceportin.IGetGoodsInfoUseCase
	getGoodsQuantityUseCase         serviceportin.IGetGoodsQuantityUseCase
	getWarehouseInfoUseCase         serviceportin.IGetWarehousesUseCase
	setMultipleGoodsQuantityUseCase serviceportin.ISetMultipleGoodsQuantityUseCase
	updateGoodDataUseCase           serviceportin.IUpdateGoodDataUseCase
}

type CatalogControllerParams struct {
	fx.In
	GetGoodsInfoUseCase             serviceportin.IGetGoodsInfoUseCase
	GetGoodsQuantityUseCase         serviceportin.IGetGoodsQuantityUseCase
	GetWarehouseInfoUseCase         serviceportin.IGetWarehousesUseCase
	SetMultipleGoodsQuantityUseCase serviceportin.ISetMultipleGoodsQuantityUseCase
	UpdateGoodDataUseCase           serviceportin.IUpdateGoodDataUseCase
}

func NewCatalogController(p CatalogControllerParams) *catalogController {
	return &catalogController{getGoodsInfoUseCase: p.GetGoodsInfoUseCase, getGoodsQuantityUseCase: p.GetGoodsQuantityUseCase, getWarehouseInfoUseCase: p.GetWarehouseInfoUseCase, setMultipleGoodsQuantityUseCase: p.SetMultipleGoodsQuantityUseCase, updateGoodDataUseCase: p.UpdateGoodDataUseCase}
}

func (cc *catalogController) getGoodsRequest(ctx context.Context, msg *nats.Msg) error { //GetGoodsInfo
	request := &request.GetGoodsInfoDTO{}

	err := json.Unmarshal(msg.Data, request)

	if err != nil {
		return nil
	}

	responseFromService := cc.getGoodsInfoUseCase.GetGoodsInfo(servicecmd.NewGetGoodsInfoCmd())

	responseToReply := response.GetGoodsDataResponseDTO{GoodMap: responseFromService.GetMap(), Err: ""}

	data, err := json.Marshal(responseToReply)

	if err != nil {
		responseToReply = response.GetGoodsDataResponseDTO{GoodMap: make(map[string]catalogCommon.Good), Err: "Cannot complete the request"}
		data, _ = json.Marshal(responseToReply)
	}

	err = msg.Respond(data)

	return err
}

func (cc *catalogController) getWarehouseRequest(ctx context.Context, msg *nats.Msg) error { //GetWarehouses
	request := &request.GetWarehousesInfoDTO{}

	err := json.Unmarshal(msg.Data, request)

	if err != nil {
		return nil
	}

	responseFromService := cc.getWarehouseInfoUseCase.GetWarehouses(servicecmd.NewGetWarehousesCmd())

	responseToReply := response.GetWarehouseResponseDTO{WarehouseMap: responseFromService.GetWarehouseMap(), Err: ""}

	data, err := json.Marshal(responseToReply)

	if err != nil {
		responseToReply = response.GetWarehouseResponseDTO{WarehouseMap: make(map[string]catalogCommon.Warehouse), Err: "Cannot complete the request"}
		data, _ = json.Marshal(responseToReply)
	}

	err = msg.Respond(data)

	return err
}

func (cc *catalogController) getGoodsGlobalQuantityRequest(ctx context.Context, msg *nats.Msg) error { //GetGoodsQuantity
	request := &request.GetGoodsQuantityDTO{}

	err := json.Unmarshal(msg.Data, request)

	if err != nil {
		return nil
	}

	responseFromService := cc.getGoodsQuantityUseCase.GetGoodsQuantity(servicecmd.NewGetGoodsQuantityCmd())

	responseToReply := response.GetGoodsQuantityResponseDTO{GoodMap: responseFromService.GetMap(), Err: ""}

	data, err := json.Marshal(responseToReply)

	if err != nil {
		responseToReply = response.GetGoodsQuantityResponseDTO{GoodMap: make(map[string]int64), Err: "Cannot complete the request"}
		data, _ = json.Marshal(responseToReply)
	}

	err = msg.Respond(data)

	return err
}

func (cc *catalogController) checkSetGoodDataRequest(request *stream.GoodUpdateData) error {
	if request.GoodID == "" || request.GoodNewName == "" || request.GoodNewDescription == "" {
		return catalogCommon.ErrRequestNotValid
	}
	return nil
}

func (cc *catalogController) setGoodDataRequest(ctx context.Context, msg jetstream.Msg) error { //AddOrChangeGoodData

	request := &stream.GoodUpdateData{}

	err := json.Unmarshal(msg.Data(), request)

	if err != nil {
		return nil
	}

	err = cc.checkSetGoodDataRequest(request)

	if err != nil {
		return err
	}

	responseFromService := cc.updateGoodDataUseCase.AddOrChangeGoodData(servicecmd.NewAddChangeGoodCmd(request.GoodID, request.GoodNewName, request.GoodNewDescription))

	if responseFromService.GetOperationResult() == catalogCommon.ErrGenericFailure {
		return catalogCommon.ErrGenericFailure
	}

	return nil
}

func (cc *catalogController) checkSetGoodQuantityRequest(request *stream.StockUpdate) error {
	if request.WarehouseID == "" || len(request.Goods) == 0 || request.Goods == nil {
		return catalogCommon.ErrRequestNotValid
	}
	return nil
}

func (cc *catalogController) setGoodQuantityRequest(ctx context.Context, msg jetstream.Msg) error { //SetMultipleGoodsQuantity

	request := &stream.StockUpdate{}

	err := json.Unmarshal(msg.Data(), request)

	if err != nil {
		return nil
	}

	err = cc.checkSetGoodQuantityRequest(request)

	if err != nil {
		return err
	}

	responseFromService := cc.setMultipleGoodsQuantityUseCase.SetMultipleGoodsQuantity(servicecmd.NewSetMultipleGoodsQuantityCmd(request.WarehouseID, request.Goods))

	if responseFromService.GetOperationResult() == catalogCommon.ErrGenericFailure {
		return catalogCommon.ErrGenericFailure
	}

	return nil
}
