package portout

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/business/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/business/response"
)

type IStockRepository interface {
	SaveStockUpdate(cmd *servicecmd.AddStockUpdateCmd) *serviceresponse.AddStockUpdateResponse
}
