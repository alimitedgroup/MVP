package port

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type ISendOrderUpdatePort interface {
	SendOrderUpdate(context.Context, SendOrderUpdateCmd) (model.Order, error)
}

type SendOrderUpdateCmd struct {
	ID           string
	Status       string
	Name         string
	Email        string
	Address      string
	CreationTime int64
	Goods        []SendOrderUpdateGood
	Reservations []string
}

type SendOrderUpdateGood struct {
	GoodId   string
	Quantity int64
}
