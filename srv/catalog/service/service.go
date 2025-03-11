package service

import (
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceportout "github.com/alimitedgroup/MVP/srv/catalog/service/portout"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
	"go.uber.org/fx"
)

type CatalogService struct {
	addOrChangeGoodDataPort serviceportout.IAddOrChangeGoodDataPort
	setGoodQuantityPort     serviceportout.ISetGoodQuantityPort
	getGoodsQuantityPort    serviceportout.IGetGoodsQuantityPort
	getGoodsInfoPort        serviceportout.IGetGoodsInfoPort
	getWarehousesPort       serviceportout.IGetWarehousesInfoPort
}

type CatalogServiceParams struct {
	fx.In
	AddOrChangeGoodDataPort serviceportout.IAddOrChangeGoodDataPort
	SetGoodQuantityPort     serviceportout.ISetGoodQuantityPort
	GetGoodsQuantityPort    serviceportout.IGetGoodsQuantityPort
	GetGoodsInfoPort        serviceportout.IGetGoodsInfoPort
	GetWarehousesPort       serviceportout.IGetWarehousesInfoPort
}

func NewCatalogService(p CatalogServiceParams) *CatalogService {
	return &CatalogService{addOrChangeGoodDataPort: p.AddOrChangeGoodDataPort, setGoodQuantityPort: p.SetGoodQuantityPort, getGoodsQuantityPort: p.GetGoodsQuantityPort, getGoodsInfoPort: p.GetGoodsInfoPort, getWarehousesPort: p.GetWarehousesPort}
}

func (cs *CatalogService) AddOrChangeGoodData(agc *servicecmd.AddChangeGoodCmd) *serviceresponse.AddOrChangeResponse {
	return cs.addOrChangeGoodDataPort.AddOrChangeGoodData(agc)
}

func checkErrinSlice(errorSlice []error) []int {
	result := []int{}
	for i := range errorSlice {
		if errorSlice[i] != nil {
			result = append(result, i)
		}
	}
	return result
}

func (cs *CatalogService) SetMultipleGoodsQuantity(cmd *servicecmd.SetMultipleGoodsQuantityCmd) *serviceresponse.SetMultipleGoodsQuantityResponse {
	warehouseID := cmd.GetWarehouseID()
	goodsSlice := cmd.GetGoods()
	var errorSlice []error
	var err error
	for i := range goodsSlice {
		err = cs.setGoodQuantityPort.SetGoodQuantity(servicecmd.NewSetGoodQuantityCmd(warehouseID, goodsSlice[i].GoodID, goodsSlice[i].Quantity)).GetOperationResult()
		errorSlice = append(errorSlice, err)
	}

	errors := checkErrinSlice(errorSlice)

	if len(errors) == 0 {
		return serviceresponse.NewSetMultipleGoodsQuantityResponse(nil, []string{})
	}

	var wrongID []string
	for i := range errors {
		wrongID = append(wrongID, goodsSlice[i].GoodID)
	}
	return serviceresponse.NewSetMultipleGoodsQuantityResponse(catalogCommon.ErrGenericFailure, wrongID)
}

func (cs *CatalogService) GetGoodsQuantity(ggqc *servicecmd.GetGoodsQuantityCmd) *serviceresponse.GetGoodsQuantityResponse {
	return cs.getGoodsQuantityPort.GetGoodsQuantity(ggqc)
}

func (cs *CatalogService) GetGoodsInfo(ggqc *servicecmd.GetGoodsInfoCmd) *serviceresponse.GetGoodsInfoResponse {
	return cs.getGoodsInfoPort.GetGoodsInfo(ggqc)
}

func (cs *CatalogService) GetWarehouses(gwc *servicecmd.GetWarehousesCmd) *serviceresponse.GetWarehousesResponse {
	return cs.getWarehousesPort.GetWarehouses(gwc)
}
