package service

import (
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
)

type IService interface {
	AddOrChangeGoodData(agc *service_Cmd.AddChangeGoodCmd) *service_Response.AddOrChangeResponse
	SetMultipleGoodsQuantity(cmd *service_Cmd.SetMultipleGoodsQuantityCmd) *service_Response.SetMultipleGoodsQuantityResponse
	GetGoodsQuantity(ggqc *service_Cmd.GetGoodsQuantityCmd)
	GetGoodsInfo(ggqc *service_Cmd.GetGoodsInfoCmd)
	GetWarehouses(gwc *service_Cmd.GetWarehousesCmd) *service_Response.GetWarehousesResponse
}
