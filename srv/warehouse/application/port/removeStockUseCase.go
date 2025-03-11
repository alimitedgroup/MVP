package port

import "context"

type IRemoveStockUseCase interface {
	RemoveStock(context.Context, RemoveStockCmd) error
}

type RemoveStockCmd struct {
	ID       string
	Quantity int64
}
