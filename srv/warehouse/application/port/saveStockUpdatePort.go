package port

import "github.com/alimitedgroup/MVP/srv/warehouse/model"

type SaveUpdateStockPort interface {
	SaveUpdateStock([]model.GoodStock) error
}
