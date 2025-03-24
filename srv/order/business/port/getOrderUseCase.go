package port

import (
	"context"
	"errors"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IGetOrderUseCase interface {
	GetOrder(context.Context, GetOrderCmd) (model.Order, error)
	GetAllOrders(context.Context) []model.Order
}

type GetOrderCmd string

var ErrOrderNotFound = errors.New("order not found")
