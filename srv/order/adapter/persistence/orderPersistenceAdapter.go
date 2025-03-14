package persistence

import (
	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type OrderPersistanceAdapter struct {
	orderRepo IOrderRepository
}

func NewOrderPersistanceAdapter(orderRepo IOrderRepository) *OrderPersistanceAdapter {
	return &OrderPersistanceAdapter{orderRepo}
}

func (s *OrderPersistanceAdapter) ApplyOrderUpdate(cmd port.ApplyOrderUpdateCmd) error {
	goods := make([]OrderUpdateGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, OrderUpdateGood{
			GoodID:   string(good.ID),
			Quantity: good.Quantity,
		})
	}

	order := Order{
		ID:           cmd.Id,
		Status:       cmd.Status,
		Name:         cmd.Name,
		Email:        cmd.Email,
		Address:      cmd.Address,
		Goods:        goods,
		CreationTime: cmd.CreationTime,
	}

	s.orderRepo.SetOrder(cmd.Id, order)
	return nil
}

func (s *OrderPersistanceAdapter) GetOrder(orderId model.OrderID) (model.Order, error) {
	order, err := s.orderRepo.GetOrder(string(orderId))
	if err != nil {
		return model.Order{}, err
	}

	modelOrder := repoOrderToModelOrder(order)
	return modelOrder, nil
}

func repoOrderToModelOrder(order Order) model.Order {
	goods := make([]model.GoodStock, 0, len(order.Goods))
	for _, good := range order.Goods {
		goods = append(goods, model.GoodStock{
			ID:       model.GoodId(good.GoodID),
			Quantity: good.Quantity,
		})
	}

	return model.Order{
		Id:           model.OrderID(order.ID),
		Name:         order.Name,
		Status:       order.Status,
		Email:        order.Email,
		Address:      order.Address,
		Goods:        goods,
		CreationTime: order.CreationTime,
	}
}
