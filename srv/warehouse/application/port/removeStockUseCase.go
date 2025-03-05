package port

import "context"

type RemoveStockUseCase interface {
	RemoveStock(context.Context, RemoveStockCmd) error
}

type RemoveStockCmd struct {
	ID       string
	Quantity int64
}
