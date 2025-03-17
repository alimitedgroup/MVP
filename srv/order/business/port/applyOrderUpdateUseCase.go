package port

import "context"

type IApplyOrderUpdateUseCase interface {
	ApplyOrderUpdate(context.Context, OrderUpdateCmd) error
}

type OrderUpdateCmd struct {
	ID           string
	Goods        []OrderUpdateGood
	Reservations []string
	Status       string
	Name         string
	FullName     string
	Address      string
	CreationTime int64
	UpdateTime   int64
}

type OrderUpdateGood struct {
	GoodID   string
	Quantity int64
}
