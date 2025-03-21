package port

import "context"

type IConfirmOrderUseCase interface {
	ConfirmOrder(context.Context, ConfirmOrderCmd) error
}

type ConfirmOrderCmd struct {
	OrderID      string
	Status       string
	Goods        []OrderUpdateGood
	Reservations []string
}

type OrderUpdateGood struct {
	GoodID   string
	Quantity int64
}
