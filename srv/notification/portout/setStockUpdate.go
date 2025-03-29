package portout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/service/response"
)

type IStockRepository interface {
	SaveStockUpdate(cmd *servicecmd.AddStockUpdateCmd) *serviceresponse.AddStockUpdateResponse
}
