package portout

import (
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/types"
)

type StockRepository interface {
	SaveStockUpdate(cmd *serviceresponse.AddStockUpdateCmd) error
}
