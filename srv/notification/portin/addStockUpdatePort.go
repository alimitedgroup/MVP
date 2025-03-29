package portin

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/business/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/business/response"
)

type StockUpdates interface {
	AddStockUpdate(cmd *servicecmd.AddStockUpdateCmd) (*serviceresponse.AddStockUpdateResponse, error)
}
