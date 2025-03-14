package port

import "github.com/alimitedgroup/MVP/srv/warehouse/model"

type IGetStockPort interface {
	GetStock(goodId model.GoodId) model.GoodStock
}
