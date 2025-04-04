package portout

import (
	"github.com/alimitedgroup/MVP/srv/notification/types"
)

type StockRepository interface {
	SaveStockUpdate(cmd *types.AddStockUpdateCmd) error
}
