package persistence

import "errors"

type IOrderRepository interface {
	GetOrder(orderId string) (Order, error)
	SetOrder(orderId string, order Order) bool
}

type Order struct {
	ID           string
	Status       string
	Name         string
	Email        string
	Address      string
	Goods        []OrderUpdateGood
	CreationTime int64
	UpdateTime   int64
}

type OrderUpdateGood struct {
	GoodID   string
	Quantity int64
}

var ErrOrderNotFound = errors.New("order not found")
