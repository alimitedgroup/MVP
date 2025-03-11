package port

import "context"

type IAddStockUseCase interface {
	AddStock(context.Context, AddStockCmd) error
}

type AddStockCmd struct {
	ID       string
	Quantity int64
}
