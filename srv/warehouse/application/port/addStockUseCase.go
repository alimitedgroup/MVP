package port

import "context"

type AddStockUseCase interface {
	AddStock(context.Context, AddStockCmd) error
}

type AddStockCmd struct {
	ID       string
	Quantity int64
}
