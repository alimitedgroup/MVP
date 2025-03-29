package portin

import (
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/types"
)

type StockUpdates interface {
	RecordStockUpdate(cmd *serviceresponse.AddStockUpdateCmd) error
}
