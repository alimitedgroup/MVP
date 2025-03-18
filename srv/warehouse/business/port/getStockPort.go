package port

import "github.com/alimitedgroup/MVP/srv/warehouse/business/model"

type IGetStockPort interface {
	GetStock(goodId model.GoodID) model.GoodStock
	GetFreeStock(goodId model.GoodID) model.GoodStock
}
