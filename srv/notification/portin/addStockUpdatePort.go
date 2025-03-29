package portin

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/service/response"
)

type IAddStockUpdateUseCase interface {
	AddStockUpdate(cmd *servicecmd.AddStockUpdateCmd) (*serviceresponse.AddStockUpdateResponse, error)
}
