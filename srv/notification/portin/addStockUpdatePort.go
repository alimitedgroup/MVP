package portin

import (
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/types"
)

type StockUpdates interface {
	AddStockUpdate(cmd *serviceresponse.AddStockUpdateCmd) (*serviceresponse.AddStockUpdateResponse, error)
}
