package port

import "context"

type ISaveOrderUpdatePort interface {
	SaveOrderUpdate(context.Context, SaveOrderUpdateCmd) error
}

type SaveOrderUpdateCmd struct {
	ID      string
	Status  string
	Name    string
	Email   string
	Address string
	Goods   []SaveOrderUpdateGood
}

type SaveOrderUpdateGood struct {
	GoodId   string
	Quantity int64
}
