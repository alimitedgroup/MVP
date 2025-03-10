package serviceportin

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
)

type IGetGoodsInfoUseCase interface {
	GetGoodsInfo(ggqc *servicecmd.GetGoodsInfoCmd) *serviceresponse.GetGoodsInfoResponse
}
