package port

import "github.com/alimitedgroup/MVP/srv/order/business/model"

type ISetCompletedWarehouseOrderPort interface {
	SetCompletedWarehouse(SetCompletedWarehouseCmd) (model.Order, error)
}

type SetCompletedWarehouseCmd struct {
	OrderId     model.OrderID
	WarehouseId string
	Goods       []model.GoodStock
}
