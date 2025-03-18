package port

import (
	"errors"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IGetStockPort interface {
	GetStock(GetStockCmd) (model.GoodStock, error)
	GetGlobalStock(GoodID model.GoodID) model.GoodStock
	GetWarehouses() []model.Warehouse
}

type GetStockCmd struct {
	WarehouseID model.WarehouseID
	GoodID      model.GoodID
}

var ErrStockNotFound = errors.New("stock not found")
