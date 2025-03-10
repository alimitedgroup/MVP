package serviceportin

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
)

type IGetGoodsQuantityUseCase interface {
	GetGoodsQuantity(ggqc *servicecmd.GetGoodsQuantityCmd) *serviceresponse.GetGoodsQuantityResponse
}
