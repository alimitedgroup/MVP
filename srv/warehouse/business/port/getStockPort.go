package port

import "github.com/alimitedgroup/MVP/srv/warehouse/business/model"

type IGetStockPort interface {
	GetStock(goodId model.GoodId) model.GoodStock
	GetFreeStock(goodId model.GoodId) model.GoodStock
}
