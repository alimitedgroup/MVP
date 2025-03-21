package port

import (
	"context"
	"errors"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IGetOrderUseCase interface {
	GetOrder(context.Context, string) (model.Order, error)
	GetAllOrders(context.Context) []model.Order
}

var ErrOrderNotFound = errors.New("order not found")
