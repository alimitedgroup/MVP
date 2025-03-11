package serviceportout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
)

type IGetWarehousesInfoPort interface {
	GetWarehouses(*servicecmd.GetWarehousesCmd) *serviceresponse.GetWarehousesResponse
}
