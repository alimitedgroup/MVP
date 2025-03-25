package port

import (
	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IApplyStockUpdatePort interface {
	ApplyStockUpdate(ApplyStockUpdateCmd)
}

type ApplyStockUpdateCmd struct {
	WarehouseID string
	Goods       []model.GoodStock
}
