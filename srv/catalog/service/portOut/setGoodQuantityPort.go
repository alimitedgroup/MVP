package service_portOut

import (
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
)

type ISetGoodQuantityPort interface {
	SetGoodQuantity(agqc *service_Cmd.SetGoodQuantityCmd) *service_Response.SetGoodQuantityResponse
}
