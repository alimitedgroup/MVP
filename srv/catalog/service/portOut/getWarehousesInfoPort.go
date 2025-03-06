package service_portOut

import (
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
)

type IGetWarehousesInfoPort interface {
	GetWarehouses(*service_Cmd.GetWarehousesCmd) *service_Response.GetWarehousesResponse
}
