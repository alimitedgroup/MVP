package persistence

import (
	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type OrderPersistanceAdapter struct {
}

func NewOrderPersistanceAdapter() *OrderPersistanceAdapter {
	return &OrderPersistanceAdapter{}
}

func (s *OrderPersistanceAdapter) GetOrder(cmd model.OrderID) (model.Order, error) {

	return model.Order{
		Id: "mock",
	}, nil
}
