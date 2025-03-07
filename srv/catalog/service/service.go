package service

import (
	"fmt"

	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
	service_portOut "github.com/alimitedgroup/MVP/srv/catalog/service/portOut"
)

type CatalogService struct {
	addOrChangeGoodDataPort service_portOut.IAddOrChangeGoodDataPort
	setGoodQuantityPort     service_portOut.ISetGoodQuantityPort
	getGoodsQuantityPort    service_portOut.IGetGoodsQuantityPort
	getGoodsInfoPort        service_portOut.IGetGoodsInfoPort
	getWarehousesPort       service_portOut.IGetWarehousesInfoPort
}

func NewCatalogService(AddOrChangeGoodDataPort service_portOut.IAddOrChangeGoodDataPort, SetGoodQuantityPort service_portOut.ISetGoodQuantityPort, GetGoodsQuantityPort service_portOut.IGetGoodsQuantityPort, GetGoodsInfoPort service_portOut.IGetGoodsInfoPort, getWarehousesPort service_portOut.IGetWarehousesInfoPort) *CatalogService {
	return &CatalogService{addOrChangeGoodDataPort: AddOrChangeGoodDataPort, setGoodQuantityPort: SetGoodQuantityPort, getGoodsQuantityPort: GetGoodsQuantityPort, getGoodsInfoPort: GetGoodsInfoPort, getWarehousesPort: getWarehousesPort}
}

func (cs *CatalogService) AddOrChangeGoodData(agc *service_Cmd.AddChangeGoodCmd) *service_Response.AddOrChangeResponse {
	return cs.addOrChangeGoodDataPort.AddOrChangeGoodData(agc)
}

func checkErrinSlice(errorSlice []string) []int {
	result := []int{}
	for i := range errorSlice {
		if errorSlice[i] != "Success" {
			result = append(result, i)
		}
	}
	return result
}

func (cs *CatalogService) SetMultipleGoodsQuantity(cmd *service_Cmd.SetMultipleGoodsQuantityCmd) *service_Response.SetMultipleGoodsQuantityResponse {
	warehouseID := cmd.GetWarehouseID()
	goodsSlice := cmd.GetGoods()
	var errorSlice []string
	var err string
	for i := range goodsSlice {
		err = cs.setGoodQuantityPort.SetGoodQuantity(service_Cmd.NewSetGoodQuantityCmd(warehouseID, goodsSlice[i].GoodID, goodsSlice[i].Quantity)).GetOperationResult()
		errorSlice = append(errorSlice, err)
	}

	errors := checkErrinSlice(errorSlice)

	if len(errors) == 0 {
		return service_Response.NewSetMultipleGoodsQuantityResponse("Success", []string{})
	}

	var wrongID []string
	for i := range errors {
		wrongID = append(wrongID, goodsSlice[i].GoodID)
	}
	fmt.Println("LOL ", errorSlice)
	return service_Response.NewSetMultipleGoodsQuantityResponse("Errors", wrongID)
}

func (cs *CatalogService) GetGoodsQuantity(ggqc *service_Cmd.GetGoodsQuantityCmd) *service_Response.GetGoodsQuantityResponse {
	return cs.getGoodsQuantityPort.GetGoodsQuantity(ggqc)
}

func (cs *CatalogService) GetGoodsInfo(ggqc *service_Cmd.GetGoodsInfoCmd) *service_Response.GetGoodsInfoResponse {
	return cs.getGoodsInfoPort.GetGoodsInfo(ggqc)
}

func (cs *CatalogService) GetWarehouses(gwc *service_Cmd.GetWarehousesCmd) *service_Response.GetWarehousesResponse {
	return cs.getWarehousesPort.GetWarehouses(gwc)
}
