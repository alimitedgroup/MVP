package service_portIn

import (
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
)

type IGetGoodsQuantityUseCase interface {
	GetGoodsQuantity(ggqc *service_Cmd.GetGoodsQuantityCmd) *service_Response.GetGoodsQuantityResponse
}
