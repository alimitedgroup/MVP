package port

import "github.com/alimitedgroup/MVP/srv/order/business/model"

type ISetCompletedWarehouseOrderPort interface {
	SetCompletedWarehouse(SetCompletedWarehouseCmd) (model.Order, error)
	SetComplete(model.OrderID) error
}

type SetCompletedWarehouseCmd struct {
	OrderId     string
	WarehouseId string
	Goods       []model.GoodStock
}
