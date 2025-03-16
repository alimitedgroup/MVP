package port

import (
	"errors"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IGetStockPort interface {
	GetStock(GetStockCmd) (model.GoodStock, error)
	GetGlobalStock(GoodID model.GoodId) model.GoodStock
	GetWarehouses() []model.Warehouse
}

type GetStockCmd struct {
	WarehouseID string
	GoodID      model.GoodId
}

var ErrStockNotFound = errors.New("stock not found")
