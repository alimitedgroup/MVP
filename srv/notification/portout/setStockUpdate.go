package portout

import (
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/types"
)

type IStockRepository interface {
	SaveStockUpdate(cmd *serviceresponse.AddStockUpdateCmd) error
}
