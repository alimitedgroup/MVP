package persistence

import (
	"errors"
)

type IOrderRepository interface {
	GetOrder(orderId string) (Order, error)
	GetOrders() []Order
	SetOrder(orderId string, order Order) bool
	AddCompletedWarehouse(orderId string, warehouseId string, goods map[string]int64) (Order, error)
	SetComplete(orderId string) error
}

type Order struct {
	ID           string
	Status       string
	Name         string
	FullName     string
	Address      string
	Goods        []OrderUpdateGood
	Warehouses   []OrderWarehouseUsed
	Reservations []string
	CreationTime int64
	UpdateTime   int64
}

type OrderUpdateGood struct {
	GoodID   string
	Quantity int64
}

type OrderWarehouseUsed struct {
	WarehouseID string
	Goods       map[string]int64
}

var ErrOrderNotFound = errors.New("order not found")
