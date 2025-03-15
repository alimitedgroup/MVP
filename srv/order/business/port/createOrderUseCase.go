package port

import (
	"context"
)

type ICreateOrderUseCase interface {
	CreateOrder(context.Context, CreateOrderCmd) (CreateOrderResponse, error)
}

type CreateOrderCmd struct {
	Name    string
	Status  string
	Email   string
	Address string
	Goods   []CreateOrderGood
}

type CreateOrderGood struct {
	GoodID   string
	Quantity int64
}

type CreateOrderResponse struct {
	OrderID string
}
