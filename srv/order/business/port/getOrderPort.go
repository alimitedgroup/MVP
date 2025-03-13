package port

import "github.com/alimitedgroup/MVP/srv/order/business/model"

type IGetOrderPort interface {
	GetOrder(model.OrderID) (model.Order, error)
}
