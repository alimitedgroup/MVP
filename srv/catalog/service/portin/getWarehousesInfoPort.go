package serviceportin

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
)

type IGetWarehousesUseCase interface {
	GetWarehouses(gwc *servicecmd.GetWarehousesCmd) *serviceresponse.GetWarehousesResponse
}
