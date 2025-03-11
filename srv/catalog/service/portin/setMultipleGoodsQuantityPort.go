package serviceportin

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
)

type ISetMultipleGoodsQuantityUseCase interface {
	SetMultipleGoodsQuantity(cmd *servicecmd.SetMultipleGoodsQuantityCmd) *serviceresponse.SetMultipleGoodsQuantityResponse
}
