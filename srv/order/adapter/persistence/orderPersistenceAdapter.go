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

func (s *OrderPersistanceAdapter) SetCompletedWarehouse(cmd port.SetCompletedWarehouseCmd) (model.Order, error) {
	goods := make(map[string]int64)
	for _, good := range cmd.Goods {
		prev, exist := goods[string(good.ID)]
		if !exist {
			prev = 0
		}
		goods[string(good.ID)] = prev + good.Quantity
	}

	order, err := s.orderRepo.AddCompletedWarehouse(string(cmd.OrderId), cmd.WarehouseId, goods)
	if err != nil {
		return model.Order{}, err
	}

	return repoOrderToModelOrder(order), nil
}

func (s *OrderPersistanceAdapter) SetComplete(orderId model.OrderID) error {
	err := s.orderRepo.SetComplete(string(orderId))
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderPersistanceAdapter) ApplyOrderUpdate(cmd port.ApplyOrderUpdateCmd) {
	var warehouses []OrderWarehouseUsed
	if old, err := s.orderRepo.GetOrder(string(cmd.Id)); err == nil {
		warehouses = old.Warehouses
	} else {
		warehouses = []OrderWarehouseUsed{}
	}

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
		FullName:     cmd.FullName,
		Address:      cmd.Address,
		Goods:        goods,
		Reservations: cmd.Reservations,
		Warehouses:   warehouses,
		UpdateTime:   cmd.UpdateTime,
		CreationTime: cmd.CreationTime,
	}

	s.orderRepo.SetOrder(cmd.Id, order)
}

func (s *OrderPersistanceAdapter) GetOrder(orderId model.OrderID) (model.Order, error) {
	order, err := s.orderRepo.GetOrder(string(orderId))
	if err != nil {
		return model.Order{}, err
	}

	modelOrder := repoOrderToModelOrder(order)
	return modelOrder, nil
}

func (s *OrderPersistanceAdapter) GetAllOrder() []model.Order {
	orders := s.orderRepo.GetOrders()

	modelOrder := repoOrdersToModelOrders(orders)
	return modelOrder
}

func repoOrderToModelOrder(order Order) model.Order {
	goods := make([]model.GoodStock, 0, len(order.Goods))
	for _, good := range order.Goods {
		goods = append(goods, model.GoodStock{
			ID:       model.GoodID(good.GoodID),
			Quantity: good.Quantity,
		})
	}

	warehouses := make([]model.OrderWarehouseUsed, 0, len(order.Warehouses))

	for _, warehouse := range order.Warehouses {
		goods := make(map[model.GoodID]int64)
		for goodId, quantity := range warehouse.Goods {
			goods[model.GoodID(goodId)] = quantity
		}
		warehouses = append(warehouses, model.OrderWarehouseUsed{
			WarehouseID: warehouse.WarehouseID,
			Goods:       goods,
		})
	}

	return model.Order{
		Id:           model.OrderID(order.ID),
		Name:         order.Name,
		Status:       order.Status,
		FullName:     order.FullName,
		Address:      order.Address,
		Goods:        goods,
		Warehouses:   warehouses,
		Reservations: order.Reservations,
		UpdateTime:   order.UpdateTime,
		CreationTime: order.CreationTime,
	}
}

func repoOrdersToModelOrders(orders []Order) []model.Order {
	list := make([]model.Order, 0, len(orders))
	for _, order := range orders {
		list = append(list, repoOrderToModelOrder(order))
	}

	return list
}
