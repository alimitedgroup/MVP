package service_portIn

import (
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
)

type IGetGoodsInfoUseCase interface {
	GetGoodsInfo(ggqc *service_Cmd.GetGoodsInfoCmd) *service_Response.GetGoodsInfoResponse
}
