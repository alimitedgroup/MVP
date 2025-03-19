package persistence

import (
	"sync"
)

type OrderRepositoryImpl struct {
	m        sync.Mutex
	orderMap map[string]Order
}

func NewOrderRepositoryImpl() *OrderRepositoryImpl {
	return &OrderRepositoryImpl{orderMap: make(map[string]Order)}
}

func (s *OrderRepositoryImpl) GetOrder(orderId string) (Order, error) {
	s.m.Lock()
	defer s.m.Unlock()

	order, exist := s.orderMap[orderId]
	if !exist {
		return Order{}, ErrOrderNotFound
	}

	return order, nil
}

func (s *OrderRepositoryImpl) GetOrders() []Order {
	s.m.Lock()
	defer s.m.Unlock()

	orders := make([]Order, 0, len(s.orderMap))
	for _, order := range s.orderMap {
		orders = append(orders, order)
	}

	return orders
}

func (s *OrderRepositoryImpl) SetOrder(orderId string, order Order) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, exist := s.orderMap[orderId]
	s.orderMap[orderId] = order

	return exist
}

func (s *OrderRepositoryImpl) AddCompletedWarehouse(orderId string, warehouseId string, goods map[string]int64) (Order, error) {
	s.m.Lock()
	defer s.m.Unlock()

	order, exist := s.orderMap[orderId]
	if !exist {
		return Order{}, ErrOrderNotFound
	}

	warehouseUsed := OrderWarehouseUsed{
		WarehouseID: warehouseId,
		Goods:       goods,
	}
	order.Warehouses = append(order.Warehouses, warehouseUsed)

	s.orderMap[orderId] = order

	return order, nil
}

func (s *OrderRepositoryImpl) SetComplete(orderId string) error {
	s.m.Lock()
	defer s.m.Unlock()

	order, exist := s.orderMap[orderId]
	if !exist {
		return ErrOrderNotFound
	}

	order.Status = "Completed"
	s.orderMap[orderId] = order

	return nil
}
