package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_portIn "github.com/alimitedgroup/MVP/srv/catalog/service/portIn"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type catalogController struct {
	getGoodsInfoUseCase             service_portIn.IGetGoodsInfoUseCase
	getGoodsQuantityUseCase         service_portIn.IGetGoodsQuantityUseCase
	getWarehouseInfoUseCase         service_portIn.IGetWarehousesUseCase
	setMultipleGoodsQuantityUseCase service_portIn.ISetMultipleGoodsQuantityUseCase
	updateGoodDataUseCase           service_portIn.IUpdateGoodDataUseCase
}

func NewCatalogController(getGoodsInfoUseCase service_portIn.IGetGoodsInfoUseCase, getGoodsQuantityUseCase service_portIn.IGetGoodsQuantityUseCase, getWarehouseInfoUseCase service_portIn.IGetWarehousesUseCase, setMultipleGoodsQuantityUseCase service_portIn.ISetMultipleGoodsQuantityUseCase, updateGoodDataUseCase service_portIn.IUpdateGoodDataUseCase) *catalogController {
	return &catalogController{getGoodsInfoUseCase: getGoodsInfoUseCase, getGoodsQuantityUseCase: getGoodsQuantityUseCase, getWarehouseInfoUseCase: getWarehouseInfoUseCase, setMultipleGoodsQuantityUseCase: setMultipleGoodsQuantityUseCase, updateGoodDataUseCase: updateGoodDataUseCase}
}

func (cc *catalogController) getGoodsRequest(ctx context.Context, msg *nats.Msg) error { //GetGoodsInfo
	request := &request.GetGoodsInfoDTO{}

	err := json.Unmarshal(msg.Data, request)

	if err != nil {
		return nil
	}

	responseFromService := cc.getGoodsInfoUseCase.GetGoodsInfo(service_Cmd.NewGetGoodsInfoCmd())

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

	responseFromService := cc.getWarehouseInfoUseCase.GetWarehouses(service_Cmd.NewGetWarehousesCmd())

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

	responseFromService := cc.getGoodsQuantityUseCase.GetGoodsQuantity(service_Cmd.NewGetGoodsQuantityCmd())

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
		return catalogCommon.NewCustomError("Not a valid request")
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

	responseFromService := cc.updateGoodDataUseCase.AddOrChangeGoodData(service_Cmd.NewAddChangeGoodCmd(request.GoodID, request.GoodNewName, request.GoodNewDescription))

	if responseFromService.GetOperationResult() == "Errors" {
		return catalogCommon.NewCustomError("An error occured")
	}

	return nil
}

func (cc *catalogController) checkSetGoodQuantityRequest(request *stream.StockUpdate) error {
	if request.WarehouseID == "" || len(request.Goods) == 0 || request.Goods == nil {
		return catalogCommon.NewCustomError("Not a valid request")
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

	responseFromService := cc.setMultipleGoodsQuantityUseCase.SetMultipleGoodsQuantity(service_Cmd.NewSetMultipleGoodsQuantityCmd(request.WarehouseID, request.Goods))

	if responseFromService.GetOperationResult() == "Errors" {
		return catalogCommon.NewCustomError("An error occured")
	}

	return nil
}
