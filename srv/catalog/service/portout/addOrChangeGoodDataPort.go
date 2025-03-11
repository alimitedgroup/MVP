package serviceportout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
)

type IAddOrChangeGoodDataPort interface {
	AddOrChangeGoodData(agc *servicecmd.AddChangeGoodCmd) *serviceresponse.AddOrChangeResponse
}
