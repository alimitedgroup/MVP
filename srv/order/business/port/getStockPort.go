package port

import (
	"errors"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IGetStockPort interface {
	GetStock(GetStockCmd) (model.GoodStock, error)
}

type GetStockCmd struct {
	WarehouseID string
	GoodID      model.GoodId
}

var ErrStockNotFound = errors.New("stock not found")
