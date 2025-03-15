package port

import "context"

type ISendOrderUpdatePort interface {
	SendOrderUpdate(context.Context, SendOrderUpdateCmd) error
}

type SendOrderUpdateCmd struct {
	ID           string
	Status       string
	Name         string
	Email        string
	Address      string
	CreationTime int64
	Goods        []SendOrderUpdateGood
}

type SendOrderUpdateGood struct {
	GoodId   string
	Quantity int64
}
