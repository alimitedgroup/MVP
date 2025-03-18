package port

import "context"

type IAddStockUseCase interface {
	AddStock(context.Context, AddStockCmd) error
}

type AddStockCmd struct {
	GoodID   string
	Quantity int64
}
