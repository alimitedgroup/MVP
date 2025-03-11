package serviceportout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
)

type ISetGoodQuantityPort interface {
	SetGoodQuantity(agqc *servicecmd.SetGoodQuantityCmd) *serviceresponse.SetGoodQuantityResponse
}
