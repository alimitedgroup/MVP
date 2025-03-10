package service

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
)

type IService interface {
	AddOrChangeGoodData(agc *servicecmd.AddChangeGoodCmd) *serviceresponse.AddOrChangeResponse
	SetMultipleGoodsQuantity(cmd *servicecmd.SetMultipleGoodsQuantityCmd) *serviceresponse.SetMultipleGoodsQuantityResponse
	GetGoodsQuantity(ggqc *servicecmd.GetGoodsQuantityCmd) *serviceresponse.GetGoodsQuantityResponse
	GetGoodsInfo(ggqc *servicecmd.GetGoodsInfoCmd) *serviceresponse.GetGoodsInfoResponse
	GetWarehouses(gwc *servicecmd.GetWarehousesCmd) *serviceresponse.GetWarehousesResponse
}
