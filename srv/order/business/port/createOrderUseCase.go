package port

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type ICreateOrderUseCase interface {
	CreateOrder(context.Context, CreateOrderCmd) (CreateOrderResponse, error)
	GetOrder(context.Context) (model.Order, error)
}

type CreateOrderCmd struct {
	Name    string
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
