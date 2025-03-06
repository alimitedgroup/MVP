package service_portIn

import (
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
)

type ISetMultipleGoodsQuantityUseCase interface {
	SetMultipleGoodsQuantity(cmd *service_Cmd.MultipleGoodsQuantityCmd) *service_Response.SetMultipleGoodsQuantityResponse
}
